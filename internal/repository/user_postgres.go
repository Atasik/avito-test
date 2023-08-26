package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	DELETION = "удаление"
	ADDITION = "добавление"
)

type User struct {
	Id int `db:"id" json:"id"`
}

type UserRepo interface {
	UpsertUser(id int, segmentsToAdd, segmentsToDelete []Segment) error
	GetUserSegments(id int) ([]Segment, error)
}

type UserPostgresqlRepo struct {
	DB *sqlx.DB
}

func NewUserPostgresqlRepo(db *sqlx.DB) *UserPostgresqlRepo {
	return &UserPostgresqlRepo{DB: db}
}

func (repo *UserPostgresqlRepo) UpsertUser(id int, segmentsToAdd, segmentsToDelete []Segment) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return ParsePostgresError(err)
	}

	query := fmt.Sprintf("INSERT INTO %s (id) VALUES($1) ON CONFLICT(id) DO NOTHING", usersTable)
	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()
		return ParsePostgresError(err)
	}

	getSegmentIdQuery := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", segmentsTable)
	insertInHistoryQuery := fmt.Sprintf("INSERT INTO %s(user_id, segment, operation, created_at) VALUES($1, $2, $3, $4)", historyTable)

	query = fmt.Sprintf("INSERT INTO %s (user_id, seg_id) VALUES($1, $2)", usersSegmentsTable)

	for _, seg := range segmentsToAdd {
		var segID int
		row := tx.QueryRow(getSegmentIdQuery, seg.Name)
		if err = row.Scan(&segID); err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}

		if _, err = tx.Exec(query, id, segID); err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}

		if _, err = tx.Exec(insertInHistoryQuery, id, seg.Name, ADDITION, time.Now()); err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND seg_id = $2", usersSegmentsTable)
	for _, seg := range segmentsToDelete {
		var segID int
		row := tx.QueryRow(getSegmentIdQuery, seg.Name)
		if err = row.Scan(&segID); err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}

		res, err := tx.Exec(query, id, segID)
		if err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}

		isDeleted, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return ParsePostgresError(err)
		}

		if isDeleted > 0 {
			if _, err = tx.Exec(insertInHistoryQuery, id, seg.Name, DELETION, time.Now()); err != nil {
				tx.Rollback()
				return ParsePostgresError(err)
			}
		}
	}

	return ParsePostgresError(tx.Commit())
}

func (repo *UserPostgresqlRepo) GetUserSegments(id int) ([]Segment, error) {
	var segments []Segment

	query := fmt.Sprintf(`SELECT s.name FROM %s s
						  INNER JOIN %s us on us.seg_id = s.id
						  INNER JOIN %s u on us.user_id = u.id
						  WHERE u.id = $1`, segmentsTable, usersSegmentsTable, usersTable)

	if err := repo.DB.Select(&segments, query, id); err != nil {
		return []Segment{}, ParsePostgresError(err)
	}
	return segments, nil
}
