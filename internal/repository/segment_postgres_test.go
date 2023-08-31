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

func TestSegmentPostgres_CreateSegment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewSegmentPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.Segment
		want    int
		wantErr bool
	}{{
		name: "OK",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", segmentsTable)).
				WithArgs("AVITO_TEST").WillReturnRows(rows)
		},
		input: domain.Segment{
			Name: "AVITO_TEST",
		},
		want: 1,
	}, {
		name: "Empty Input",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id"})
			mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", segmentsTable)).
				WithArgs("").WillReturnRows(rows)
		},
		input: domain.Segment{
			Name: "",
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input)
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

func TestSegmentPostgres_DeleteSegment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewSegmentPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.Segment
		wantErr bool
	}{{
		name: "OK",
		mock: func() {
			mock.ExpectExec(fmt.Sprintf("DELETE FROM %s WHERE (.+)", segmentsTable)).
				WithArgs("AVITO_TEST").WillReturnResult(sqlmock.NewResult(0, 1))
		},
		input: domain.Segment{
			Name: "AVITO_TEST",
		},
	},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s WHERE (.+)", segmentsTable)).
					WithArgs("AVITO_TEST").WillReturnError(sql.ErrNoRows)
			},
			input: domain.Segment{
				Name: "AVITO_TEST",
			},
			wantErr: true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
