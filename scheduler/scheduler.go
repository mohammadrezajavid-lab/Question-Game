package scheduler

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
	"sync"
	"time"
)

type Config struct {
	MatchingJobCrontab string `mapstructure:"matching_job_crontab"`
	QuizPoolJobCrontab string `mapstructure:"quizpool_job_crontab"`
}

type Scheduler struct {
	sch         gocron.Scheduler
	matchingSvc matchingservice.Service
	quizSvc     quizservice.Service
	config      Config
}

func New(matchingSvc matchingservice.Service, quizSvc quizservice.Service, config Config) Scheduler {

	sch, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		metrics.FailedCreateSchCounter.Inc()
		logger.Fatal(err, errormessage.ErrorMsgFailedCreateSch)
	}

	return Scheduler{
		sch:         sch,
		matchingSvc: matchingSvc,
		quizSvc:     quizSvc,
		config:      config,
	}
}

// Start Start() is a long running process
func (s *Scheduler) Start(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	metrics.StartSchCounter.Inc()
	logger.Info(infomessage.InfoMsgSchStart)

	s.newJobMatchWaitedUser()
	s.newJobGenerateQuiz()
	s.sch.Start()

	<-ctx.Done()

	logger.Info(infomessage.InfoMsgSchExiting)

	if sErr := s.sch.Shutdown(); sErr != nil {
		logger.Warn(sErr, errormessage.ErrorMsgShutdownSch)
	}
}
