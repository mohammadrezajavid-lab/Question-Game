package scheduler

import (
	"fmt"
	"log"
	"time"
)

type Scheduler struct {
}

func New() *Scheduler {
	return &Scheduler{}
}

// Start Start() is a long running process
func (s *Scheduler) Start(done <-chan bool) {
	log.Println("scheduler started.")

	for {
		select {
		case d := <-done:
			fmt.Printf("\nScheduler exiting %v...\n", d)

			return
		default:
			fmt.Println("implement me ")
			time.Sleep(time.Millisecond * 2000)
		}

	}
}
