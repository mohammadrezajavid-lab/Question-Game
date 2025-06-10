package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
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

	sch, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		metrics.FailedCreateSchCounter.Inc()
		logger.Fatal(err, errormessage.ErrorMsgFailedCreateSch)
	}

	return Scheduler{
		sch:         sch,
		matchingSvc: matchingSvc,
		config:      config,
	}
}

// Start Start() is a long running process
func (s *Scheduler) Start(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	metrics.StartSchCounter.Inc()
	logger.Info(infomessage.InfoMsgSchStart)

	s.newJobMatchWaitedUser()
	s.sch.Start()

	<-ctx.Done()

	logger.Info(infomessage.InfoMsgSchExiting)

	if sErr := s.sch.Shutdown(); sErr != nil {
		logger.Warn(sErr, errormessage.ErrorMsgShutdownSch)
	}
}

func (s *Scheduler) newJobMatchWaitedUser() {
	matchWaitedUsersJob, nErr := s.sch.NewJob(
		gocron.CronJob(s.config.Crontab, false),
		gocron.NewTask(s.matchWaitedUserTask),
		gocron.WithSingletonMode(gocron.LimitModeWait),
		gocron.WithName("match-waited-user"),
		gocron.WithTags("matching-service"),
	)

	if nErr != nil {
		metrics.FailedCreateMatchWaitedUserJobCounter.Inc()
		logger.Fatal(nErr, errormessage.ErrorMsgFailedStartMatchWaitedUserJob)
	}

	matchWaitedUsersJobInfo := fmt.Sprintf("job info: name[%s], uuid[%v], tags[%v]",
		matchWaitedUsersJob.Name(), matchWaitedUsersJob.ID(), matchWaitedUsersJob.Tags())

	metrics.CreateMatchWaitedUserJobCounter.Inc()
	logger.Info(fmt.Sprintf("matchWaitedUser job created, %s", matchWaitedUsersJobInfo))
}

func (s *Scheduler) matchWaitedUserTask() {

	ctx, cancel := context.WithTimeout(context.Background(), s.matchingSvc.GetConfig().ContextTimeOut)
	defer cancel()
	err := s.matchingSvc.MatchWaitedUsers(ctx)
	if err != nil {
		metrics.MatchWaitedUserFailedJobCounter.Inc()
		logger.Info(fmt.Sprintf("matchWaitedUsers_failed, %v", err))
	} else {
		metrics.MatchWaitedUserRunSuccessfullyJobCounter.Inc()
		logger.Info("matchWaitedUsers_successfully")
	}
}
