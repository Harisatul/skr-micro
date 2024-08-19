package tryout

import "time"

const (
	orderCancellationDuration      = 15 * time.Second
	CreateOrderTopic               = "tryout.create_order"
	CancelOrderTopic               = "tryout.cancel_order"
	CompleteOrderTopic             = "tryout.complete_order"
	NotifyCancelOrderTopic         = "scheduler.notify_cancel_order"
	NotifyCancelOrderConsumerGroup = "tryout.notify_cancel_order"
)
