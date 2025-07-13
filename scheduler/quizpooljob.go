package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"time"
)

func (s *Scheduler) newJobGenerateQuiz() {
	generateQuizJob, err := s.sch.NewJob(
		gocron.CronJob(s.config.QuizPoolJobCrontab, false),
		gocron.NewTask(s.generateQuizTask),
		gocron.WithSingletonMode(gocron.LimitModeWait),
		gocron.WithName("generate-quiz"),
		gocron.WithTags("quizpool-service"),
	)

	if err != nil {
		metrics.FailedCreateGenerateQuizJobCounter.Inc()
		logger.Fatal(err, errormessage.ErrorMsgFailedStartGenerateQuizJob)
	}

	generateQuizJobInfo := fmt.Sprintf("job info: name[%s], uuid[%v], tags[%v]",
		generateQuizJob.Name(), generateQuizJob.ID(), generateQuizJob.Tags())

	metrics.CreateMatchWaitedUserJobCounter.Inc()
	logger.Info(fmt.Sprintf("generateQuiz job created, %s", generateQuizJobInfo))
}

func (s *Scheduler) generateQuizTask() {
	logger.Info(fmt.Sprintf("generateQuizTask started at: %s", time.Now().Format(time.RFC3339)))
	s.quizSvc.GenerateQuiz()

	metrics.GenerateQuizRunSuccessfullyJobCounter.Inc()
	logger.Info("generateQuizJob_successfully")
}
