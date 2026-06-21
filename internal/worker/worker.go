package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/ink-pwd/WatchKeeper/internal/queue"
	"github.com/ink-pwd/WatchKeeper/internal/scheduler"
	"github.com/ink-pwd/WatchKeeper/internal/telegram"
	"github.com/ink-pwd/WatchKeeper/internal/utils"
)

type Worker struct {
	queue *queue.Queue
	sched *scheduler.Scheduler
	tg    *telegram.TelegramBot
}

func NewWorker(queue *queue.Queue, sched *scheduler.Scheduler, tg *telegram.TelegramBot) *Worker {
	return &Worker{
		queue: queue,
		sched: sched,
		tg:    tg,
	}
}

func (w *Worker) Start() {
	var (
		task   *queue.Task
		ctx    context.Context
		cancel context.CancelFunc
		alive  bool
		msg    string
	)
	for task = range w.queue.Channel() {
		/*
			Ограничиваем время получения ответа на 3 секунды.
			Что бы не создавать бесконечные подключения.
		*/
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		alive = utils.IsServerAlive(ctx, task.Addr)
		cancel()
		if alive {
			/*
				Если сервер жив, добавляем его в конец очереди для следующей проверки
			*/
			w.sched.AddTask(task.ChatID, task.Addr)
		} else {
			msg = fmt.Sprintf("🔴 \"%s\" doesn't answer. It has been removed from the review list",
				task.Addr)
			w.tg.SendMessage(task.ChatID, msg)
		}
	}
}
