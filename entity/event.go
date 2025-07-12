package entity

type Event string

const (
	MatchingUsersMatchedEvent = "matching.matched_users" // "serviceName.eventName"
	GameSvcCreatedGameEvent   = "game.game_created"
	GameSvcGameEvent          = "game.game_event"
)
