package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	sch         gocron.Scheduler
	matchingSvc *matchingservice.Service
}

func New(matchingSvc *matchingservice.Service) *Scheduler {

	// TODO - we can set location timezone in config.yaml file and use this for scheduler
	sch, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	return &Scheduler{
		sch:         sch,
		matchingSvc: matchingSvc,
	}
}

// Start Start() is a long running process
func (s *Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	log.Println("scheduler started.")

	s.newJobMatchWaitedUser()
	s.sch.Start()

	<-done
	fmt.Printf("Scheduler exiting...\n")
	wg.Done()
}

func (s *Scheduler) newJobMatchWaitedUser() {
	matchWaitedUsersJob, nErr := s.sch.NewJob(
		gocron.CronJob("*/2 * * * *", false),
		gocron.NewTask(s.matchWaitedUser),
		gocron.WithSingletonMode(gocron.LimitModeWait),
		gocron.WithName("match-waited-user"),
		gocron.WithTags("matching-service"),
	)

	matchWaitedUsersJobInfo := fmt.Sprintf("job info: name[%s], uuid[%v], tags[%v]",
		matchWaitedUsersJob.Name(), matchWaitedUsersJob.ID(), matchWaitedUsersJob.Tags())

	if nErr != nil {
		log.Fatalf("Failed to create matchWaitedUser job, %s\n", matchWaitedUsersJobInfo)
	}

	log.Printf("matchWaitedUser job created, %s\n", matchWaitedUsersJobInfo)
}

func (s *Scheduler) matchWaitedUser() {
	log.Println("Executing matchWaitedUser job at:", time.Now())

	s.matchingSvc.MatchWaitedUser(param.NewMatchWaitedUserRequest())
}
