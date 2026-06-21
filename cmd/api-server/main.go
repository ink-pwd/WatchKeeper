package main

import (
	"context"
	"log"
	"time"

	"github.com/ink-pwd/WatchKeeper/internal/config"
	"github.com/ink-pwd/WatchKeeper/internal/queue"
	"github.com/ink-pwd/WatchKeeper/internal/scheduler"
	"github.com/ink-pwd/WatchKeeper/internal/telegram"
	"github.com/ink-pwd/WatchKeeper/internal/worker"
	"github.com/redis/go-redis/v9"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var (
		cfg   *config.Config
		bot   *tgbotapi.BotAPI
		tg    *telegram.TelegramBot
		rdb   *redis.Client
		qu    *queue.Queue
		sched *scheduler.Scheduler
		w     *worker.Worker
		i     int
		err   error
	)
	cfg, err = config.GetConfig()
	if err != nil {
		log.Fatalf("[Error] Get config: %s", err.Error())
	}

	bot, err = tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatalf("[Error] Connect to bot: %s", err.Error())
	}
	log.Printf("[Info] Authorized on account: %s", bot.Self.UserName)

	rdb = redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[Error] Connect to redis: %s", err.Error())
	}
	log.Printf("[Info] Redis is connected to \"%s\"", cfg.RedisAddr)

	/*
		Создаем новую очередь и планировщик.
		Один планировщик запускаем асинхронно.
	*/
	qu = queue.NewQueue(cfg.BufferQueue)

	sched = scheduler.NewScheduler(rdb, time.Duration(cfg.Interval)*time.Second, qu)

	go sched.Run()

	/*
		Создаем бота и обработчик запросов из очереди
	*/
	tg = telegram.CreateTelegramBot(bot, sched, cfg.TimeOutTelegram)
	w = worker.NewWorker(qu, sched, tg)
	/*
		Асинхронно запускаем worker`ов
	*/
	for i = range cfg.WorkerPoolSize {
		log.Printf("[Info] %d worker start", i+1)
		go w.Start()
	}
	tg.ListenServ()
}
