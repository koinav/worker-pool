package main

import (
	"fmt"
	"time"
	"worker-pool/pkg/workerpool"
)

// Пример использования
func main() {
	pool := workerpool.NewWorkerPool(3, 100)

	// Отправка нескольких задач
	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		if err := pool.Submit(task); err != nil {
			fmt.Printf("%s: Ошибка при отправке задачи: %s\n", time.Now().Format(time.RFC3339), err)
		}
	}

	time.Sleep(time.Second * 2)

	// Удаление воркера
	if err := pool.RemoveWorker(); err != nil {
		fmt.Printf("%s: Ошибка при удалении воркера: %s\n", time.Now().Format(time.RFC3339), err)
	}

	time.Sleep(time.Second * 2)

	// Добавление нового воркера
	if err := pool.AddWorker(); err != nil {
		fmt.Printf("%s: Ошибка при добавлении воркера: %s\n", time.Now().Format(time.RFC3339), err)
	}

	time.Sleep(time.Second * 5)

	// Принудительное завершение всех воркеров
	pool.StopAll()

	if err := pool.Submit("task after stop"); err != nil {
		fmt.Printf("%s: Ошибка при отправке задачи: %s\n", time.Now().Format(time.RFC3339), err)
	}
}
