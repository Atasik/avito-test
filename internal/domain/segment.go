package domain

import (
	"segmenter/pkg/timejson"
)

type Segment struct {
	ID         int                 `db:"id" swaggerignore:"true" json:"-"`
	Name       string              `db:"name" json:"name" example:"AVITO_TEST" validate:"required"`
	Percentage float32             `db:"percentage" json:"percentage,omitempty" example:"0.25"`
	ExpiredAt  *timejson.CivilTime `db:"expired_at" json:"expired_at,omitempty" swaggertype:"primitive,string" example:"2024-03-23"`
}
