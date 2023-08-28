package domain

import "segmenter/pkg/timejson"

type History struct {
	UserID    int                `db:"user_id"`
	Segment   string             `db:"segment"`
	Operation string             `db:"operation"`
	CreatedAt timejson.CivilTime `db:"created_at"`
}
