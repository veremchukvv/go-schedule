package main

import (
	"flag"
	"sync"
	"time"

	"github.com/golang/glog"
)

type Sched struct {
	Fn func()
	wg sync.WaitGroup
}

func NewSched(fn func()) *Sched {
	return &Sched{Fn: fn}
}

func main() {
	flag.Parse()

	sched := NewSched(fn)
	workCh := make(chan *Sched, 1)
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

			for lastRun := range workCh { // Пытаемся дочитать из закрытого канала
				glog.Info("Last run execution...")
				lastRun.Fn() // Если что-то осталось, то выполняем функцию ещё раз
			}
			break LOOP
		case <-ticker.C:
			workCh <- sched
		case f := <-workCh:
			glog.Info("starting func execution...")
			f.wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer f.wg.Done()
				f.Fn() // Запускаем нашу функцию внутри горутины
			}(&f.wg)
			f.wg.Wait()
		}
	}
}

func fn() {
	time.Sleep(time.Millisecond * 1000) // Имитация 1 секунды длительности работы
	glog.Info(">>fn")                   // Вывод в лог факта завершения выполнения
}
