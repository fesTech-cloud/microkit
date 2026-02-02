package messaging

import "errors"

var (
	ErrPublishFailed   = errors.New("failed to publish message")
	ErrSubscriptionErr = errors.New("failed to subscribe to topic")
)
