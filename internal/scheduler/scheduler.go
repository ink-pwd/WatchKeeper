package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ink-pwd/WatchKeeper/internal/queue"
	"github.com/redis/go-redis/v9"
)

/*
Сохраняем и автоматически сортируем запросы на мониторинг используя ZAdd Redis.
При получении запроса даем сигнал в канал, что бы тот начал обработку.
Первый запрос всегда первый на выполнение, получается своеобразная FIFO очередь.
Обрабатываем все запросы поочереди, так как нет необходимости в асинхроне.
Далее, дождавшись время на выполнение, передаем запрос в очередь на выполнение.
*/

const DefaultQueueKey = "monitor:queue"

type TaskMember struct {
	ChatID int64  `json:"chat_id"`
	Addr   string `json:"addr"`
}

type Scheduler struct {
	rdb      *redis.Client
	interval time.Duration
	signalCh chan struct{}
	queue    *queue.Queue
}

func NewScheduler(rdb *redis.Client, interval time.Duration,
	queue *queue.Queue) *Scheduler {
	return &Scheduler{
		rdb:      rdb,
		interval: interval,
		signalCh: make(chan struct{}, 1),
		queue:    queue,
	}
}

func (s *Scheduler) AddTask(chat_id int64, addr string) error {
	var (
		task []byte
		err  error
	)

	task, err = json.Marshal(TaskMember{
		ChatID: chat_id,
		Addr:   addr,
	})

	if err != nil {
		return err
	}

	err = s.rdb.ZAdd(context.Background(), DefaultQueueKey, redis.Z{
		Score:  float64(time.Now().Add(s.interval).Unix()),
		Member: string(task),
	}).Err()

	if err != nil {
		return err
	}

	/*
		Даем сигнал планировщику, что задача добавлена.
	*/
	s.Signal()

	return nil
}

/*
Получаем сигнал, а это значит, что уже как минимум одна задача добавлена.
Так как при выполнении у нас time.Sleep мы можем пропустить новые сигналы о добавлении.
Поэтому мы находимся в taskProcessing и получаем новые задачи пока они есть не ожидая сигнал.
*/
func (s *Scheduler) Run() {
	for range s.signalCh {
		s.taskProcessing()
	}
}

func (s *Scheduler) taskProcessing() {
	var (
		result    []redis.Z
		err       error
		nextCheck int64
		now       int64
		data      []byte
		task      TaskMember
	)
	for {
		/*
			Берем элемент удаляя его из очереди.
		*/
		result, err = s.rdb.ZPopMin(context.Background(), DefaultQueueKey, 1).Result()
		if err != nil || len(result) == 0 {
			/*
				Выходим из цикла если не получили элемент или получили ошибку
			*/
			return
		}

		/*
			Получаем информацию о задаче.
		*/
		data = []byte(result[0].Member.(string))
		err = json.Unmarshal(data, &task)
		if err != nil {
			/*
				Если unmarshall выдал ошибку, значит переходим на следующую итерацию цикла.
				Так как могли быть добавлены новые запросы пока мы обрабатывали этот.
			*/
			continue
		}
		nextCheck = int64(result[0].Score)
		now = time.Now().Unix()
		if nextCheck > now {
			/*
				Если время на выполнение не пришло - time sleep.
			*/
			time.Sleep(time.Duration(nextCheck-now) * time.Second)
		}
		/*
			Добавляем задачу в очередь на выполнение.
		*/
		s.queue.Push(&queue.Task{
			ChatID: task.ChatID,
			Addr:   task.Addr,
		})
	}
}

func (s *Scheduler) Signal() {
	select {
	case s.signalCh <- struct{}{}:
	default:
	}
}
