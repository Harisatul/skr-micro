package tryout

import "time"

const (
	orderCancellationDuration      = 15 * time.Second
	CreateSubmissionTopic          = "exam.create_submission"
	CancelOrderTopic               = "tryout.cancel_order"
	CreateAnswerTopic              = "exam.create_answer"
	NotifyCancelOrderTopic         = "scheduler.notify_cancel_order"
	NotifyCancelOrderConsumerGroup = "tryout.notify_cancel_order"
)
