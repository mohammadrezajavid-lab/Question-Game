package quizparam

import "golang.project/go-fundamentals/gameapp/entity"

type GetQuizRequest struct {
	Category   entity.Category
	Difficulty entity.QuestionDifficulty
}

type GetQuizResponse struct {
	QuestionIds []uint
}
