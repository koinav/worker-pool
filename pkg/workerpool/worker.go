package workerpool

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker представляет отдельного воркера
type worker struct {
	id       int
	quit     chan struct{}
	tasks    <-chan string
	wg       *sync.WaitGroup
	stopOnce sync.Once
}

// newWorker создает нового воркера
func newWorker(id int, tasks <-chan string, wg *sync.WaitGroup) *worker {
	return &worker{
		id:    id,
		quit:  make(chan struct{}),
		tasks: tasks,
		wg:    wg,
	}
}

// start запускает воркер для обработки задач
func (w *worker) start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			// Проверяем, не поступил ли сигнал стоп
			select {
			case <-w.quit:
				fmt.Printf("%s: Worker %d: завершение работы по сигналу стоп.\n", time.Now().Format(time.RFC3339), w.id)
				return
			default:
				// перходим к обработке задач, не блокируя горутину
			}

			select {
			case task, ok := <-w.tasks:
				if !ok {
					// Если канал закрыт, завершаем работу
					fmt.Printf("%s: Worker %d: входной канал закрыт. Завершение работы.\n", time.Now().Format(time.RFC3339), w.id)
					return
				}
				// Обработка задачи
				processingTime := time.Duration(rand.Intn(10)) * time.Second
				fmt.Printf("%s: Worker %d: начал обработку %s\n", time.Now().Format(time.RFC3339), w.id, task)
				time.Sleep(processingTime)
				fmt.Printf("%s: Worker %d: обработал данные %s (время выполнения - %.0f сек.)\n",
					time.Now().Format(time.RFC3339),
					w.id,
					task,
					processingTime.Seconds())
			default:
				// если задач нет, немного ждем и повторяем действия
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
}

// stop завершает работу воркера
func (w *worker) stop() {
	fmt.Printf("%s: Worker %d: поступил сигнал стоп.\n", time.Now().Format(time.RFC3339), w.id)
	w.stopOnce.Do(func() {
		close(w.quit)
	})
}
