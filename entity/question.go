package entity

type PossibleAnswerChoice uint8

const (
	A PossibleAnswerChoice = iota + 1
	B
	C
	D
)

func (p PossibleAnswerChoice) IsValid() bool {

	if p >= A && p <= D {

		return true
	}

	return false
}

type PossibleAnswer struct {
	Id     uint
	Text   string
	Choice PossibleAnswerChoice
}

type QuestionDifficulty uint8

const (
	Easy QuestionDifficulty = iota + 1
	Medium
	Hard
)

func (q QuestionDifficulty) IsValid() bool {

	if q >= Easy && q <= Hard {

		return true
	}

	return false
}

type Question struct {
	Id              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswer   uint // id possibleAnswer
	Difficulty      QuestionDifficulty
	CategoryId      uint
}

func NewQuestion(id uint, text string) *Question {
	return &Question{
		Id:              id,
		Text:            text,
		PossibleAnswers: nil,
		CorrectAnswer:   0,
		Difficulty:      0,
		CategoryId:      0,
	}
}
