package leaderboard

import "errors"

var (
	errTicketNotFound = errors.New("ticket: not found")
	errOrderExist     = errors.New("leaderboard: already exist")
	errOrderNotFound  = errors.New("leaderboard: not found")
)
