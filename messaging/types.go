package messaging

import "sync"

type Broker struct {
	mu     sync.Mutex
	topics map[string]Topic
}

type Topic struct {
	subscribers []Subscriber
	queue       []Message
}

type Message struct {
	ID           string
	Body         string
	PUBLISHED_AT string
	CONSUMED_AT  string
}

type Subscriber struct {
	ID            string
	Endpoint      string
	Topics        []string
	LastHeartbeat string
	status        SubscriberStatus
}

type SubscriberStatus string

const (
	ACTIVE   SubscriberStatus = "active"
	INACTIVE SubscriberStatus = "inactive"
	FAILED   SubscriberStatus = "failed"
)

type SubscribeRequest struct {
	Endpoint string   `json:"endpoint"`
	Topics   []string `json:"topics"`
}

type UnsubscribeRequest struct {
	SubscriberID string   `json:"subscriber_id"`
	Topics       []string `json:"topics"`
}

type PublishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type RetrieveRequest struct {
	SubscriberId string `json:"sub_id"`
	Topic        string `json:"topic"`
	StartDate    string `json:"startDate"`
}
