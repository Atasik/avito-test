package domain

type User struct {
	ID int `db:"id" json:"id" example:"1" validate:"required,gt=0"`
}
