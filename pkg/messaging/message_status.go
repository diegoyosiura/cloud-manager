package messaging

type MessageStatus int

const (
	NotSent   MessageStatus = 0
	Queued    MessageStatus = 1
	Sending   MessageStatus = 2
	Sent      MessageStatus = 3
	SendError MessageStatus = 4
)
