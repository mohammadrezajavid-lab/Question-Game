package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"log"
	"sync"
	"time"
)

type Config struct {
	Crontab string `mapstructure:"crontab"`
}

type Scheduler struct {
	sch         gocron.Scheduler
	matchingSvc matchingservice.Service
	config      Config
}

func New(matchingSvc matchingservice.Service, config Config) Scheduler {

	// TODO - we can set location timezone in config.yaml file and use this for scheduler
	sch, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	return Scheduler{
		sch:         sch,
		matchingSvc: matchingSvc,
		config:      config,
	}
}

// Start Start() is a long running process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	log.Println("scheduler started.")

	s.newJobMatchWaitedUser()
	s.sch.Start()

	<-done
	fmt.Printf("Scheduler exiting...\n")
	time.Sleep(15 * time.Second)
	wg.Done()
}

func (s Scheduler) newJobMatchWaitedUser() {
	matchWaitedUsersJob, nErr := s.sch.NewJob(
		gocron.CronJob(s.config.Crontab, false),
		gocron.NewTask(s.matchWaitedUserTask),
		gocron.WithSingletonMode(gocron.LimitModeWait),
		gocron.WithName("match-waited-user"),
		gocron.WithTags("matching-service"),
	)

	matchWaitedUsersJobInfo := fmt.Sprintf("job info: name[%s], uuid[%v], tags[%v]",
		matchWaitedUsersJob.Name(), matchWaitedUsersJob.ID(), matchWaitedUsersJob.Tags())

	if nErr != nil {
		log.Fatalf("Failed to create matchWaitedUser job, %s\n", matchWaitedUsersJobInfo)
	}

	log.Printf("✅ matchWaitedUser job created, %s\n", matchWaitedUsersJobInfo)
}

func (s Scheduler) matchWaitedUserTask() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.matchingSvc.MatchWaitedUsers(ctx)
	if err != nil {
		log.Printf("🚨 MatchWaitedUsers failed: %v", err)
		// TODO - update metrics
	} else {
		log.Println("✅ MatchWaitedUsers ran successfully.")
		// TODO - update metrics
	}
}
