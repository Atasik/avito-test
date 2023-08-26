package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type History struct {
	UserId    int       `db:"user_id"`
	Segment   string    `db:"segment"`
	Operation string    `db:"operation"`
	CreatedAt time.Time `db:"created_at"`
}

type HistoryRepo interface {
	GetHistoryForPeriod(start, end time.Time, userId int) ([]History, error)
}

type HistoryPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewHistoryPostgresqlRepo(db *sqlx.DB) *HistoryPostgresqlRepo {
	return &HistoryPostgresqlRepo{DB: db}
}

func (repo *HistoryPostgresqlRepo) GetHistoryForPeriod(start, end time.Time, userId int) ([]History, error) {
	var history []History

	query := fmt.Sprintf(`SELECT user_id, segment, operation, created_at FROM %s WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3`, historyTable)

	if err := repo.DB.Select(&history, query, userId, start, end); err != nil {
		return []History{}, ParsePostgresError(err)
	}
	return history, nil
}
