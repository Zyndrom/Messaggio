package kafka

import (
	"context"
	"encoding/json"
	"messaggio/internal/models"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type kafkaService struct {
	writer         *kafka.Writer
	reader         *kafka.Reader
	messageService messageService
}

type messageService interface {
	ProcessMessage(msg models.Message)
}

func New(messageService messageService) *kafkaService {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")
	address := host + ":" + port
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, "topic", 0)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    "topic",
		Balancer: &kafka.LeastBytes{},
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     "topic",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	k := &kafkaService{
		writer:         writer,
		reader:         reader,
		messageService: messageService,
	}
	go k.MessageReader()
	return k
}

func (k *kafkaService) SendNewMessage(msg models.Message) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(strconv.Itoa(msg.ID)),
			Value: message,
		},
	)

	if err != nil {
		logrus.Error("Failed to write messages:", err)
	}
	logrus.Info("Kafka: New message")
	return err
}

func (k *kafkaService) MessageReader() {

	for {
		m, err := k.reader.ReadMessage(context.Background())
		logrus.Info("Kafka: Read message")
		if err != nil {
			logrus.Error("Kafka read message error:", err)
			continue
		}
		var msg models.Message
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			logrus.Error("Kafka unmarshal message error:", err)
			continue
		}
		k.messageService.ProcessMessage(msg)
	}
}
