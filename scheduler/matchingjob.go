package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"time"
)

func (s *Scheduler) newJobMatchWaitedUser() {
	matchWaitedUsersJob, nErr := s.sch.NewJob(
		gocron.CronJob(s.config.MatchingJobCrontab, false),
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

	logger.Info(fmt.Sprintf("matchWaitedUsers started at: %s", time.Now().Format(time.RFC3339)))

	ctx, cancel := context.WithTimeout(context.Background(), s.config.MatchingJobContextTimeOut)
	defer cancel()

	s.matchingSvc.MatchWaitedUsers(ctx)

	metrics.MatchWaitedUserRunSuccessfullyJobCounter.Inc()
	logger.Info("matchWaitedUsers_successfully")
}
