package repository

import (
	"messaggio/internal/models"
	"messaggio/internal/repository/postgres"
)

type storage interface {
	SaveMessageAndSetID(msg *models.Message) error
	ProcessMessage(msg models.Message) error
	GetStatistic() (models.Statistic, error)
}

type repository struct {
	storage storage
}

func New() *repository {
	return &repository{
		storage: postgres.New(),
	}
}

func (r *repository) SaveMessageAndSetID(msg *models.Message) error {
	return r.storage.SaveMessageAndSetID(msg)
}
func (r *repository) ProcessMessage(msg models.Message) error {
	return r.storage.ProcessMessage(msg)
}

func (r *repository) GetStatistic() (models.Statistic, error) {
	return r.storage.GetStatistic()
}
