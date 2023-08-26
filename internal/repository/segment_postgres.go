package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Segment struct {
	Name string `db:"name" json:"name"`
}

type SegmentRepo interface {
	CreateSegment(seg Segment) (int, error)
	DeleteSegment(seg Segment) error
	// GetSegments() ([]Segment, error)
}

type SegmentPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewSegmentPostgresqlRepo(db *sqlx.DB) *SegmentPostgresqlRepo {
	return &SegmentPostgresqlRepo{DB: db}
}

func (repo *SegmentPostgresqlRepo) CreateSegment(seg Segment) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", segmentsTable)

	row := repo.DB.QueryRow(query, seg.Name)
	err := row.Scan(&id)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return id, nil
}

func (repo *SegmentPostgresqlRepo) DeleteSegment(seg Segment) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name = $1", segmentsTable)

	_, err := repo.DB.Exec(query, seg.Name)
	if err != nil {
		return ParsePostgresError(err)
	}
	return nil
}

// func (repo *SegmentPostgresqlRepo) GetSegments() ([]Segment, error) {
// 	var segments []Segment

// 	query := fmt.Sprintf("SELECT * FROM %s", segmentsTable)

// 	if err := repo.DB.Select(&segments, query); err != nil {
// 		return nil, ParsePostgresError(err)
// 	}

// 	return segments, nil
// }
