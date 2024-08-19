package grade

import "github.com/google/uuid"

type submissionRequest struct {
	TryoutID string `json:"tryout_id"`
	Token    string `json:"token"`
}

type answerRequest struct {
	SubmissionID string `json:"user_test_submission_id"`
	QuestionID   string `json:"question_id"`
	ChoiceID     string `json:"choice_id"`
}

type submissionResponse struct {
	Id      uuid.UUID `json:"id"`
	version uint32    `json:"version"`
}

type answerResponse struct {
	Id      uuid.UUID `json:"id"`
	version uint32    `json:"version"`
}

type cancelOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket cancelOrderTicketPayload `json:"ticket"`
}
type cancelOrderTicketPayload struct {
	Version uint32 `json:"version"`
}
type notifyCancelOrderMessage struct {
	ID uint64 `json:"id"`
}
type submissionMessage struct {
	Id       string `json:"id"`
	TryOutId string `json:"tryout_id"`
	Token    string `json:"token"`
	Version  int64  `json:"version"`
}

type gradeMessage struct {
	Id       string `json:"id"`
	Token    string `json:"token"`
	TryoutID string `json:"tryout_id"`
	Score    int    `json:"score"`
	Version  int64  `json:"version"`
}

type answerMessage struct {
	Id           string `json:"id"`
	ChoiceId     string `json:"token"`
	SubmissionID string `json:"submission_id"`
	Version      int64  `json:"version"`
}
