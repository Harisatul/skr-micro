package leaderboard

import "time"

const (
	orderCancellationDuration    = 15 * time.Second
	CreateGradeTopic             = "grade.create_grade"
	UpdateGradeTopic             = "grade.update_grade"
	GradeSubmissionConsumerGroup = "grade.submission_group"
	GradeUpdateConsumerGroup     = "grade.update_group"
)
