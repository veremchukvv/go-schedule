package main

import (
	"flag"
	"sync"
	"time"

	"github.com/golang/glog"
)

type Sched struct {
	Fn      func()
	LastRun bool
	mu      sync.Mutex
}

func NewSched(fn func()) *Sched {
	return &Sched{Fn: fn}
}

func main() {
	flag.Parse()

	sched := NewSched(fn)
	sched.Start()
	// sched.Stop()
}

func fn() {
	time.Sleep(time.Millisecond * 1000) // Имитация 1 секунды длительности работы
	glog.Info(">>fn")                   // Вывод в лог факта завершения выполнения
}

func (s *Sched) Start() {
	var wg sync.Once

	ch := time.After(time.Millisecond * 1530)

LOOP:
	for {
		select {
		case <-ch:
			glog.Info(">>exit")
			if s.LastRun {
				wg.Do(fn)
			}
			break LOOP
		default:
			if s.mu.TryLock() {
				go func() {
					fn()
					s.mu.Unlock()
				}() // Вызов функции в отдельном потоке
			} else {
				s.LastRun = true
			}
		}
	}
}

// func (s *Sched) Stop() {
// 	var wg sync.Once

// 	if s.LastRun {
// 		wg.Do(fn)
// 	}
// }
