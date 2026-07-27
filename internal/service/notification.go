package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/rajabhishekmaurya/ecom/internal/model"
	"github.com/rajabhishekmaurya/ecom/libs/config"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Broker),
			Topic:    cfg.Kafka.Topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Publish(event *model.EventPayment) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer() *KafkaConsumer {

	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "payment-events",
			GroupID: "notification-service",
		}),
	}
}

func (c *KafkaConsumer) Start() error {

	fmt.Println("Notification Service is waiting for events...")

	for {

		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event model.EventPayment

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("====================================")
		fmt.Println("Notification Received")
		fmt.Println("Order ID      :", event.OrderID)
		fmt.Println("TransactionID :", event.TransactionID)
		fmt.Println("Amount        :", event.Amount)
		fmt.Println("Status        :", event.Status)
		fmt.Println("Email Sent Successfully")
		fmt.Println("====================================")
	}
}
