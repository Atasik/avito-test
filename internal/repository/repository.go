package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	SegmentRepo
	UserRepo
	HistoryRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		SegmentRepo: NewSegmentPostgresqlRepo(db),
		UserRepo:    NewUserPostgresqlRepo(db),
		HistoryRepo: NewHistoryPostgresqlRepo(db),
	}
}
