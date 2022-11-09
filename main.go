package main

import (
	"flag"
	"sync"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	workCh := make(chan func(), 1)
	ticker := time.NewTicker(time.Millisecond) // Имитируем попытки запуска функции по таймеру
	glog.Info("Start")
	cancelCh := time.After(time.Millisecond * 1530) // Задаем таймер на работу планировщика

LOOP:
	for {
		select {
		case <-cancelCh: // Инициируем остановку планировщика по сигналу из канала
			glog.Info(">>exit")
			close(workCh) // Закрываем канал
			ticker.Stop() // Останавливаем счётчик

			for range workCh { // Пытаемся дочитать из закрытого канала
				glog.Info("Last run execution...")
				fn() // Если что-то осталось, то выполняем функцию ещё раз
			}
			break LOOP
		case <-ticker.C:
			workCh <- fn
		case f := <-workCh:
			glog.Info("starting func execution...")
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				f() // Запускаем нашу функцию внутри горутины
			}(&wg)
			wg.Wait()
		}
	}
}

func fn() {
	time.Sleep(time.Millisecond * 1000) // Имитация 1 секунды длительности работы
	glog.Info(">>fn")                   // Вывод в лог факта завершения выполнения
}
