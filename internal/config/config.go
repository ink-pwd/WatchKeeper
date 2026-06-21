package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Token           string
	RedisAddr       string
	TimeOutTelegram int
	Interval        int
	BufferQueue     int
	WorkerPoolSize  int
}

func GetConfig() (*Config, error) {
	var (
		cfg Config
		str string
		err error
	)
	_ = godotenv.Load()

	cfg.Token = os.Getenv("TOKEN")

	cfg.RedisAddr = os.Getenv("REDISADDR")

	str = os.Getenv("TIMEOUTTELEGRAM")
	cfg.TimeOutTelegram, err = strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	str = os.Getenv("INTERVAL")
	cfg.Interval, err = strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	str = os.Getenv("BUFFERQUEUE")
	cfg.BufferQueue, err = strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	str = os.Getenv("WORKERPOOLSIZE")
	cfg.WorkerPoolSize, err = strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
