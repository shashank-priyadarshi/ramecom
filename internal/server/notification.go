package server

// starter is anything that blocks in Start (e.g. Kafka consumer).
type starter interface {
	Start() error
}

// NotificationServer runs a long-lived consumer loop.
type NotificationServer struct {
	consumer starter
}

func NewNotification(consumer starter) *NotificationServer {
	return &NotificationServer{consumer: consumer}
}

func (s *NotificationServer) Start() error {
	return s.consumer.Start()
}
