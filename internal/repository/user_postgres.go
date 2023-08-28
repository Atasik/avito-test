package repository

import (
	"fmt"
	"segmenter/internal/domain"
	"segmenter/pkg/postgres"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	DELETION = "delete"
	ADDITION = "add"
)

type UserRepo interface {
	UpsertUser(userID int, segmentsToAdd, segmentsToDelete []domain.Segment) error
	GetUserSegments(userID int) ([]domain.Segment, error)
	DeleteExpiredSegments() error
}

type UserPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewUserPostgresqlRepo(db *sqlx.DB) *UserPostgresqlRepo {
	return &UserPostgresqlRepo{DB: db}
}

func (repo *UserPostgresqlRepo) UpsertUser(userID int, segmentsToAdd, segmentsToDelete []domain.Segment) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return postgres.ParsePostgresError(err)
	}

	query := fmt.Sprintf("INSERT INTO %s (id) VALUES($1) ON CONFLICT(id) DO NOTHING", usersTable)
	if _, err = tx.Exec(query, userID); err != nil {
		tx.Rollback() //nolint:errcheck
		return postgres.ParsePostgresError(err)
	}

	getSegmentIDQuery := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", segmentsTable)

	queryInterval := fmt.Sprintf("INSERT INTO %s (user_id, seg_id) VALUES($1, $2)", usersSegmentsTable)
	query = fmt.Sprintf("INSERT INTO %s (user_id, seg_id, expired_at) VALUES($1, $2, $3)", usersSegmentsTable)

	for _, seg := range segmentsToAdd {
		var segID int
		row := tx.QueryRow(getSegmentIDQuery, seg.Name)
		if err = row.Scan(&segID); err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		if time.Time(seg.ExpiredAt).IsZero() {
			if _, err = tx.Exec(queryInterval, userID, segID); err != nil {
				tx.Rollback() //nolint:errcheck
				return postgres.ParsePostgresError(err)
			}
		} else {
			if _, err = tx.Exec(query, userID, segID, time.Time(seg.ExpiredAt)); err != nil {
				tx.Rollback() //nolint:errcheck
				return postgres.ParsePostgresError(err)
			}
		}
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND seg_id = $2", usersSegmentsTable)
	for _, seg := range segmentsToDelete {
		var segID int
		row := tx.QueryRow(getSegmentIDQuery, seg.Name)
		if err = row.Scan(&segID); err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		_, err := tx.Exec(query, userID, segID)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}
	}

	return postgres.ParsePostgresError(tx.Commit())
}

func (repo *UserPostgresqlRepo) GetUserSegments(userID int) ([]domain.Segment, error) {
	var segments []domain.Segment

	query := fmt.Sprintf(`SELECT s.name FROM %s s
						  INNER JOIN %s us on us.seg_id = s.id
						  INNER JOIN %s u on us.user_id = u.id
						  WHERE u.id = $1`, segmentsTable, usersSegmentsTable, usersTable)

	if err := repo.DB.Select(&segments, query, userID); err != nil {
		return []domain.Segment{}, postgres.ParsePostgresError(err)
	}
	return segments, nil
}

func (repo *UserPostgresqlRepo) DeleteExpiredSegments() error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE expired_at <= $1`, usersSegmentsTable)

	if _, err := repo.DB.Exec(query, time.Now()); err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}
