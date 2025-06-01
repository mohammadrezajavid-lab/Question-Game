package entity

type Presence struct {
	UserId    uint
	Timestamp int64
}

func NewPresence(userId uint, timestamp int64) Presence {
	return Presence{UserId: userId, Timestamp: timestamp}
}
