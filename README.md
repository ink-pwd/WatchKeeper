# WatchKeeper

WatchKeeper is a Telegram bot that monitors website availability and notifies users when a website becomes unavailable.

The service uses Redis as a persistent scheduler storage and automatically restores monitoring tasks after application restarts.

## Features

* Monitor website availability via HTTP requests
* Telegram notifications when a website becomes unavailable
* Redis-based persistent task storage
* Automatic recovery after service restart
* Scheduler based on Redis Sorted Sets (ZSET)
* Worker Pool for concurrent website checks
* Docker support

## Architecture

WatchKeeper consists of several components:

### Telegram Bot

Handles user commands and manages monitoring subscriptions.

### Scheduler

Stores monitoring tasks in Redis Sorted Sets and determines when each website should be checked.

### Worker Pool

Processes website availability checks concurrently.

### Redis

Used as persistent storage for monitoring tasks and scheduler data.

## Project Structure

```text
WatchKeeper/
├── cmd/
│   └── api-server/
│       └── main.go
│
├── internal/
│   ├── config/
│       └── config.go
│   ├── queue/
│       └── queue.go
│   ├── scheduler/
│       └── scheduler.go
│   ├── telegram/
│       └── telegram.go
│   ├── utils/
│       ├── parse_test.go
│       └── parse.go
│   └── worker/
│       └── worker.go
│
├── .env
├── dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## How It Works

1. User adds a website through Telegram.
2. The website is placed into a Redis Sorted Set.
3. Scheduler waits until the next execution time.
4. A worker performs an HTTP request to the website.
5. If the website is unavailable, the user receives a Telegram notification.
6. The website is re-scheduled for the next check.

## Recovery After Restart

Monitoring tasks are stored in Redis.

If the application is restarted, WatchKeeper automatically restores all monitoring tasks and continues monitoring without data loss.

## Configuration

Create a `.env` file based on `.env.example`.

Required variables:

```env
TOKEN=your token tg

REDISADDR=redis:6379
TIMEOUTTELEGRAM=30

INTERVAL=300 server operation check interval in seconds

BUFFERQUEUE=100
WORKERPOOLSIZE=2
```

## Run With Docker

```bash
docker compose up --build -d
```

## Technologies

* Go
* Telegram Bot API
* Redis
* Docker
* Docker Compose

## Future Improvements

* Per-user monitoring settings
* Monitoring history
* Response time statistics
* Website status dashboard
* Multiple notification channels

## License

MIT
