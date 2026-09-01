package telegram

type Handler func(Context) error

type Route struct {
	Pattern string

	MatchType MatchType

	Handler Handler

	Middleware []Middleware
}

type MatchType int

const (
	MatchExact MatchType = iota

	MatchPrefix

	MatchContains

	MatchCommand

	MatchCommandStartOnly

	MatchCallbackData
)
