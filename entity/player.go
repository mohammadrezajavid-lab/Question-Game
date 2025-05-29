package entity

type MatchedUsers struct {
	Category Category
	UserIds  []uint
}

func NewMatchedUsers(category Category, userIds []uint) MatchedUsers {
	return MatchedUsers{
		Category: category,
		UserIds:  userIds,
	}
}
