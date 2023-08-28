package repository

import (
	"fmt"
	"segmenter/internal/domain"
	"segmenter/pkg/postgres"
	"time"

	"github.com/jmoiron/sqlx"
)

type HistoryRepo interface {
	GetHistoryForPeriod(period time.Time, userID int) ([]domain.History, error)
}

type HistoryPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewHistoryPostgresqlRepo(db *sqlx.DB) *HistoryPostgresqlRepo {
	return &HistoryPostgresqlRepo{DB: db}
}

func (repo *HistoryPostgresqlRepo) GetHistoryForPeriod(period time.Time, userID int) ([]domain.History, error) {
	var history []domain.History

	query := fmt.Sprintf(`SELECT user_id, segment, operation, created_at FROM %s WHERE user_id = $1 AND date_part('year', created_at) = $2 AND date_part('month', created_at) = $3`, historyTable)

	if err := repo.DB.Select(&history, query, userID, period.Year(), period.Month()); err != nil {
		return []domain.History{}, postgres.ParsePostgresError(err)
	}
	return history, nil
}
