package repository

import (
	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	segmentsTable      = "segments"
	usersSegmentsTable = "users_segments"
	historyTable       = "history"
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
