package postgres

import (
	"database/sql"
	"fmt"
	"messaggio/internal/models"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type postgresql struct {
	db *sql.DB
}

func New() *postgresql {
	usr := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		usr, pass, host, port, dbName)
	logrus.Info(connStr)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf(err.Error())
	}
	psg := &postgresql{db: db}
	return psg
}

func (p *postgresql) SaveMessageAndSetID(msg *models.Message) error {
	query := `INSERT INTO messages (content, processed) VALUES ($1, $2) returning id;`
	err := p.db.QueryRow(query, msg.Content, msg.Processed).Scan(&msg.ID)
	return err
}

func (p *postgresql) ProcessMessage(msg models.Message) error {
	query := `UPDATE messages SET processed_at = CURRENT_TIMESTAMP, processed = TRUE WHERE id = $1 AND processed = FALSE;`

	_, err := p.db.Exec(query, msg.ID)
	return err
}

func (p *postgresql) GetStatistic() (models.Statistic, error) {
	timeQuery := `SELECT 
	EXTRACT(EPOCH FROM AVG(processed_at - created_at)) AS average_processing_time_seconds
	FROM messages
	WHERE processed_at IS NOT NULL;`
	totalMessageQuery := `SELECT count(*) from messages where processed = TRUE;`

	var stat models.Statistic
	row := p.db.QueryRow(timeQuery)
	err := row.Scan(&stat.AverageProcessingTime)
	if err != nil {
		return stat, err
	}
	row = p.db.QueryRow(totalMessageQuery)
	err = row.Scan(&stat.ProcessedMessages)
	if err != nil {
		return stat, err
	}
	return stat, nil
}
