package gameparam

import "golang.project/go-fundamentals/gameapp/entity"

type GetQuizRequest struct {
	GameId uint `json:"game_id"`
}

type GetQuizResponse struct {
	GameId    uint       `json:"game_id"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Id              uint             `json:"id"`
	Text            string           `json:"text"`
	PossibleAnswers []PossibleAnswer `json:"possible_answers"`
	CorrectAnswer   uint             `json:"correct_answer"` // id possibleAnswer
	Difficulty      uint8            `json:"difficulty"`
	Category        string           `json:"category"`
}

type PossibleAnswer struct {
	Id     uint   `json:"id"`
	Text   string `json:"text"`
	Choice uint8  `json:"choice"`
}

func MapPossibleAnswersToParam(possibleAnswers []entity.PossibleAnswer) []PossibleAnswer {

	paramPAS := make([]PossibleAnswer, 0, len(possibleAnswers))
	paramPA := PossibleAnswer{
		Id:     0,
		Text:   "",
		Choice: 0,
	}

	for _, pa := range possibleAnswers {
		paramPA.Id = pa.Id
		paramPA.Text = pa.Text
		paramPA.Choice = uint8(pa.Choice)

		paramPAS = append(paramPAS, paramPA)
	}

	return paramPAS
}
