package grade

import "errors"

var (
	errTicketNotFound = errors.New("ticket: not found")
	errOrderExist     = errors.New("grade: already exist")
	errOrderNotFound  = errors.New("grade: not found")
)
