package entity

type MatchedUsers struct {
	Category   Category
	Difficulty QuestionDifficulty
	UserIds    []uint
}

func NewMatchedUsers(category Category, difficulty QuestionDifficulty, userIds []uint) MatchedUsers {
	return MatchedUsers{
		Category:   category,
		Difficulty: difficulty,
		UserIds:    userIds,
	}
}
