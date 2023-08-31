package repository

import (
	"database/sql"
	"fmt"
	"segmenter/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserPostgres_UpsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserPostgresqlRepo(sqlxDB)

	type args struct {
		userID                          int
		segmentsToAdd, segmentsToDelete []domain.Segment
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{{
		name: "OK Add",
		mock: func() {
			mock.ExpectBegin()

			mock.ExpectExec(fmt.Sprintf("INSERT INTO %s ", usersTable)).WithArgs(1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s WHERE name IN", segmentsTable)).WithArgs("TEST").
				WillReturnRows(rows)

			mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", usersSegmentsTable)).
				WillReturnResult(sqlmock.NewResult(0, 1))

			mock.ExpectCommit()
		},
		input: args{
			userID: 1,
			segmentsToAdd: []domain.Segment{
				{
					Name: "TEST",
				},
			},
		},
	},
		{
			name: "OK Delete",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s ", usersTable)).WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s WHERE name IN", segmentsTable)).WithArgs("TEST").
					WillReturnRows(rows)

				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s WHERE seg_id IN", usersSegmentsTable)).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			input: args{
				userID: 1,
				segmentsToDelete: []domain.Segment{
					{
						Name: "TEST",
					},
				},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s ", usersTable)).WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s WHERE name IN", segmentsTable)).WithArgs("TEST").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			input: args{
				userID: 1,
				segmentsToAdd: []domain.Segment{
					{
						Name: "TEST",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpsertUserSegments(tt.input.userID, tt.input.segmentsToAdd, tt.input.segmentsToDelete)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_GetUserSegments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.User
		want    []domain.Segment
		wantErr bool
	}{{
		name: "OK",
		mock: func() {
			rows := sqlmock.NewRows([]string{"name"}).
				AddRow("AVITO_TEST")

			mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s s INNER JOIN %s us on (.+) INNER JOIN %s u on (.+) WHERE (.+)", segmentsTable, usersSegmentsTable, usersTable)).
				WithArgs(1000).WillReturnRows(rows)
		},
		input: domain.User{
			ID: 1000,
		},
		want: []domain.Segment{
			{Name: "AVITO_TEST"},
		},
	},
		{
			name: "Not found",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s s INNER JOIN %s us on (.+) INNER JOIN %s u on (.+) WHERE (.+)", segmentsTable, usersSegmentsTable, usersTable)).
					WithArgs(404).WillReturnError(sql.ErrNoRows)
			},
			input: domain.User{
				ID: 404,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSegments(tt.input.ID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
