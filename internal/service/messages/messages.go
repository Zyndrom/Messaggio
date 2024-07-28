package messages

import (
	"messaggio/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

type storage interface {
	SaveMessageAndSetID(msg *models.Message) error
	ProcessMessage(msg models.Message) error
	GetStatistic() (models.Statistic, error)
}

type kafka interface {
	SendNewMessage(msg models.Message) error
}

type messageService struct {
	storage storage
	kafka   kafka
}

func New(storage storage) *messageService {
	return &messageService{
		storage: storage,
	}
}
func (m *messageService) SetKafka(kafka kafka) {
	m.kafka = kafka
}

func (m *messageService) NewMessage(text string) error {
	msg := models.Message{
		Content:   text,
		Processed: false,
		CreatedAt: time.Now(),
	}
	err := m.storage.SaveMessageAndSetID(&msg)
	if err != nil {
		logrus.Errorf("Failed to save message: %s, error: %v", msg.Content, err)
		return err
	}
	err = m.kafka.SendNewMessage(msg)
	if err != nil {
		logrus.Errorf("Failed send message to kafka. Message: %s, error: %v", msg.Content, err)
		return err
	}
	return nil
}

func (m *messageService) ProcessMessage(msg models.Message) {
	err := m.storage.ProcessMessage(msg)
	if err != nil {
		logrus.Error("Failed process message. Id:", msg.ID)
	}
}

func (m *messageService) GetStatistic() (models.Statistic, error) {
	return m.storage.GetStatistic()
}
