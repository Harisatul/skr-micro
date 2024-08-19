package tryout

import (
	"github.com/google/uuid"
	"time"
)

type UserTestSubmission struct {
	ID                           uuid.UUID
	CreatedAt                    int64
	UpdatedAt                    int64
	ParticipantsID               string
	Token                        string
	TryoutID                     uuid.UUID
	SubmissionStart              time.Time
	SubmissionEnd                time.Time
	UserSubmittedAnswerExercises []UserSubmittedAnswer
	Version                      int64
}

type UserSubmittedAnswer struct {
	ID                   uuid.UUID
	CreatedAt            int64
	UpdatedAt            int64
	UserTestSubmissionID uuid.UUID
	QuestionID           uuid.UUID
	ChoiceID             uuid.UUID
	Version              int64
}

type TryOut struct {
	ID          string
	CreatedAt   int64
	UpdatedAt   int64
	TryOutName  string
	TestType    string
	TryOutQuota int
	version     int64
}
