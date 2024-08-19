package grade

import "time"

const (
	orderCancellationDuration    = 15 * time.Second
	CreateSubmissionTopic        = "exam.create_submission"
	CreateAnswerTopic            = "exam.create_answer"
	CreateGradeTopic             = "grade.create_grade"
	UpdateGradeTopic             = "grade.update_grade"
	GradeSubmissionConsumerGroup = "grade.submission_group"
	GradeAnswerConsumerGroup     = "grade.answer_group"
)
