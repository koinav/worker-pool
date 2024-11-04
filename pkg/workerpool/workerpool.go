package workerpool

import (
	"fmt"
	"sync"
	"time"
)

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
	workerID  int
	workers   []*worker
	tasks     chan string
	wg        sync.WaitGroup
	mu        sync.Mutex
	closeOnce sync.Once
	isClosed  bool
}

// NewWorkerPool создает новый пул воркеров с заданным количеством воркеров
// и размером буферизированного канала
func NewWorkerPool(workers, capacity int) *WorkerPool {
	pool := &WorkerPool{
		workers: make([]*worker, 0),
		tasks:   make(chan string, capacity),
	}

	for i := 0; i < workers; i++ {
		if err := pool.AddWorker(); err != nil {
			fmt.Printf("Ошибка при добавлении воркера: %s\n", err)
		}
	}

	return pool
}

// Submit отправляет задачу во входной канал
func (p *WorkerPool) Submit(task string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isClosed {
		return fmt.Errorf("невозможно добавить задачу, пул воркеров закрыт")
	}
	p.tasks <- task
	return nil
}

// Close закрывает входной канал
// Воркеры продолжат работать до выполнения всех задач в канале
func (p *WorkerPool) Close() {
	p.closeOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.isClosed = true
		close(p.tasks)
		fmt.Printf("%s: Входной канал закрыт закрыт, добавление задач невозможно.\n", time.Now().Format(time.RFC3339))
	})
}

// StopAll завершает работу всех воркеров и закрывает входной канал
// Воркеры завершат выполняемые задачи и остановятся
func (p *WorkerPool) StopAll() {
	p.closeOnce.Do(func() {
		fmt.Printf("%s: Получен сигнал завершения.\n", time.Now().Format(time.RFC3339))
		p.mu.Lock()
		defer p.mu.Unlock()
		p.isClosed = true
		close(p.tasks)
		for _, worker := range p.workers {
			worker.stop()
		}
		p.wg.Wait()
		fmt.Printf("%s: Воркерпул закрыт, все процессы остановлены.\n", time.Now().Format(time.RFC3339))
	})
}

// AddWorker добавляет нового воркера в пул
func (p *WorkerPool) AddWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isClosed {
		return fmt.Errorf("невозможно добавить воркера: пул воркеров закрыт")
	}

	p.workerID++
	w := newWorker(p.workerID, p.tasks, &p.wg)
	p.workers = append(p.workers, w)
	fmt.Printf("%s: Worker %d: добавлен в пул.\n", time.Now().Format(time.RFC3339), w.id)
	w.start()
	return nil
}

// RemoveWorker удаляет последнего добавленного воркера из пула
// Если воркер занят выполнением задачи, он завершит ее перед остановкой
func (p *WorkerPool) RemoveWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	l := len(p.workers)
	if l == 0 {
		return fmt.Errorf("нет активных воркеров в пуле")
	}

	worker := p.workers[l-1]
	p.workers = p.workers[:l-1]
	fmt.Printf("%s: Worker %d: будет удален из пула.\n", time.Now().Format(time.RFC3339), worker.id)
	worker.stop()
	return nil
}
