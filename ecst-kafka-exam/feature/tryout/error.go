package tryout

import "errors"

var (
	errTicketNotFound = errors.New("ticket: not found")
	errOrderExist     = errors.New("tryout: already exist")
	errOrderNotFound  = errors.New("tryout: not found")
)
