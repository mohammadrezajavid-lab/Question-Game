package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	sch *gocron.Scheduler
}

func New() *Scheduler {
	return &Scheduler{}
}

// Start Start() is a long running process
func (s *Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	log.Println("scheduler started.")

	for {
		select {
		case d := <-done:
			for i := 0; i < 10; i++ {
				time.Sleep(700 * time.Millisecond)
				fmt.Println(time.Now())
			}
			fmt.Printf("\nScheduler exiting %v...\n", d)

			wg.Done()

			return
		default:
			fmt.Println("implement me ")
			time.Sleep(time.Millisecond * 200)
		}

	}
}
