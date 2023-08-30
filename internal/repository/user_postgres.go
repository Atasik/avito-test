package repository

import (
	"fmt"
	"log"
	"segmenter/internal/domain"
	"segmenter/pkg/postgres"
	"strings"
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

	selectQuery := fmt.Sprintf("SELECT id FROM %s WHERE name in (?)", segmentsTable)

	if len(segmentsToAdd) != 0 {
		var segmentsBuilder strings.Builder

		idsToAdd := make([]int, len(segmentsToAdd))
		namesToAdd := make([]string, len(segmentsToAdd))
		for i, seg := range segmentsToAdd {
			namesToAdd[i] = seg.Name
		}

		query, args, err := sqlx.In(selectQuery, namesToAdd)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		query = repo.DB.Rebind(query)

		if err = repo.DB.Select(&idsToAdd, query, args...); err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}
		segmentsBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (user_id, seg_id, expired_at) VALUES ", usersSegmentsTable))

		args = []interface{}{}
		argId := 1
		for idx, seg := range segmentsToAdd {
			args = append(args, userID)
			args = append(args, idsToAdd[idx])
			args = append(args, seg.ExpiredAt)
			segmentsBuilder.WriteString(fmt.Sprintf(`($%d,$%d,$%d),`, argId, argId+1, argId+2))
			argId += 3
		}

		query = strings.TrimSuffix(segmentsBuilder.String(), ",")

		_, err = repo.DB.Exec(query, args...)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}
	}

	if len(segmentsToDelete) != 0 {
		idsToDelete := make([]int, len(segmentsToDelete))
		namesToDelete := make([]string, len(segmentsToDelete))
		for i, seg := range segmentsToDelete {
			namesToDelete[i] = seg.Name
		}

		query, args, err := sqlx.In(selectQuery, namesToDelete)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		query = repo.DB.Rebind(query)

		if err = repo.DB.Select(&idsToDelete, query, args...); err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		query, args, err = sqlx.In(fmt.Sprintf("DELETE FROM %s WHERE seg_id IN (?) AND user_id = ? ", usersSegmentsTable), idsToDelete, userID)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}

		query = repo.DB.Rebind(query)

		_, err = repo.DB.Exec(query, args...)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return postgres.ParsePostgresError(err)
		}
	}

	return postgres.ParsePostgresError(tx.Commit())
}

func (repo *UserPostgresqlRepo) GetUserSegments(userID int) ([]domain.Segment, error) {
	var segments []domain.Segment

	query := fmt.Sprintf(`SELECT s.name, us.expired_at FROM %s s
						  INNER JOIN %s us on us.seg_id = s.id
						  INNER JOIN %s u on us.user_id = u.id
						  WHERE u.id = $1`, segmentsTable, usersSegmentsTable, usersTable)

	if err := repo.DB.Select(&segments, query, userID); err != nil {
		return []domain.Segment{}, postgres.ParsePostgresError(err)
	}

	log.Println(segments)
	return segments, nil
}

func (repo *UserPostgresqlRepo) DeleteExpiredSegments() error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE expired_at <= $1`, usersSegmentsTable)

	if _, err := repo.DB.Exec(query, time.Now()); err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}
