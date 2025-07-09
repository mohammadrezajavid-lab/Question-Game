package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
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
