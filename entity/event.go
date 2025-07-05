package entity

type Event string

const (
	MatchingUsersMatchedEvent = "matching.matched_users" // "serviceName.eventName"
	GameSvcCreatedGameEvent   = "game.created_games"
)
