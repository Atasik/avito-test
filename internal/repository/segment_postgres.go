package repository

import (
	"fmt"
	"segmenter/internal/domain"
	"segmenter/pkg/postgres"

	"github.com/jmoiron/sqlx"
)

type SegmentRepo interface {
	CreateSegment(seg domain.Segment) (int, error)
	DeleteSegment(seg domain.Segment) error
}

type SegmentPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewSegmentPostgresqlRepo(db *sqlx.DB) *SegmentPostgresqlRepo {
	return &SegmentPostgresqlRepo{DB: db}
}

func (repo *SegmentPostgresqlRepo) CreateSegment(seg domain.Segment) (int, error) {
	var id int
	var query string
	if seg.Percentage > 0 {
		query = fmt.Sprintf("INSERT INTO %s (name, percentage) VALUES ($1, $2) RETURNING id", segmentsTable)
	} else {
		query = fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", segmentsTable)
	}

	row := repo.DB.QueryRow(query, seg.Name)
	err := row.Scan(&id)
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	return id, nil
}

func (repo *SegmentPostgresqlRepo) DeleteSegment(seg domain.Segment) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name = $1", segmentsTable)

	_, err := repo.DB.Exec(query, seg.Name)
	if err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}
