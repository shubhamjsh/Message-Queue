package messaging

import (
	"Message-Queue/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]Topic),
	}

}

func NewMessage(body string) Message {
	return Message{
		ID:           uuid.New().String(),
		Body:         body,
		PUBLISHED_AT: time.Now().Format(time.RFC3339),
		CONSUMED_AT:  "",
	}
}

func NewSubscriber(endpoint string, topics []string) Subscriber {
	return Subscriber{
		ID:            uuid.New().String(),
		Endpoint:      endpoint,
		Topics:        topics,
		LastHeartbeat: time.Now().Format(time.RFC3339),
		status:        ACTIVE,
	}
}

func (b *Broker) Subscribe(endpoint string, topics []string) (string, error) {
	if !utils.IsValidURL(endpoint) {
		return "", fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}
	subscriber := NewSubscriber(endpoint, topics)
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topicName := range topics {
		if _, exists := b.topics[topicName]; !exists {
			newTopic := Topic{
				subscribers: []Subscriber{},
				queue:       []Message{},
			}
			b.topics[topicName] = newTopic
		}

		isSubscriberExists := false
		for _, sub := range b.topics[topicName].subscribers {
			if sub.Endpoint == endpoint {
				isSubscriberExists = true
				break
			}
		}
		if !isSubscriberExists {
			topic := b.topics[topicName]
			topic.subscribers = append(topic.subscribers, subscriber)
			b.topics[topicName] = topic
		}
	}
	return subscriber.ID, nil
}

func (b *Broker) Unsubscribe(subscriberID string, topics []string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topicName := range topics {
		topic, exists := b.topics[topicName]
		if !exists {
			return fmt.Errorf("topic %s not found", topicName)
		}
		var index int = 0
		for index, sub := range b.topics[topicName].subscribers {
			if sub.ID == subscriberID {
				topic.subscribers = append(topic.subscribers[:index], topic.subscribers[index+1:]...)
				break
			}
		}
		if len(topic.subscribers) == index {
			return fmt.Errorf("subscriber %s not subscribed to topic %s", subscriberID, topicName)
		}
	}
	return nil
}

func (b *Broker) Retrieve(subID string, topicName string,
	startDate string) ([]Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return nil, fmt.Errorf("topic %s not found", topicName)
	}

	var subscriberExists bool
	for _, sub := range topic.subscribers {
		if sub.ID == subID {
			subscriberExists = true
			break
		}
	}
	if !subscriberExists {
		return nil, fmt.Errorf("subscriber %s does not exist for topic %s", subID, topicName)
	}
	if startDate == "" {
		return topic.queue, nil
	}
	var messsag []Message
	for _, msg := range topic.queue {
		messsag = append(messsag, msg)
	}
	return messsag, nil
}

func (b *Broker) Publish(topicName, message string) error {
	topic, exists := b.topics[topicName]
	if !exists {
		return fmt.Errorf("topic %s not found", topicName)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	newMessage := NewMessage(message)
	topic.queue = append(topic.queue, newMessage)
	err := b.propagateMsg(&topic)
	if err != nil {
		return err
	}
	b.topics[topicName] = topic
	return nil
}

func (b *Broker) propagateMsg(topic *Topic) error {
	for _, msg := range topic.queue {
		if msg.CONSUMED_AT == "" {
			for _, sub := range topic.subscribers {
				if sub.status == ACTIVE {
					body := []byte(`{
						"Message": "` + msg.Body + `",
						"Published_at": "` + msg.PUBLISHED_AT + `",
						}`)
					request, err := http.NewRequest("POST", sub.Endpoint, bytes.NewBuffer(body))
					if err != nil {
						return fmt.Errorf("failed to construct request for subscriber%s: %v", sub.ID, err)
					}

					request.Header.Set("Content-Type", "application/json; charset=utf-8")
					client := &http.Client{}
					res, err := client.Do(request)
					if err != nil {
						return fmt.Errorf("failed to send message to subscriber %s: %v", sub.ID, err)
					}
					_, err = io.ReadAll(res.Body)
					if err != nil {
						return fmt.Errorf("failed to parse the response from subscriber %s: %v", sub.ID, err)
					}
					defer res.Body.Close()
					fmt.Println("Message sent to subscriber", sub.ID)
					sub.LastHeartbeat = time.Now().Format(time.RFC3339)
				}
			}
			msg.CONSUMED_AT = time.Now().Format(time.RFC3339)
		}
	}
	return nil
}
