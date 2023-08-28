package domain

import "segmenter/pkg/timejson"

type Segment struct {
	Name      string             `db:"name" json:"name"`
	ExpiredAt timejson.CivilTime `db:"expired_at" json:"expired_at"`
}
