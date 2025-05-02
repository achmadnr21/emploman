package domain

import (
	"time"
)

type Grade struct {
	ID         int       `json:"id" db:"id"`
	Code       string    `json:"code" db:"code"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

type GradeInterface interface {
	FindAll() ([]Grade, error)
	FindByID(id int) (*Grade, error)
	Save(grade *Grade) (*Grade, error)
	Update(grade *Grade) (*Grade, error)
	Delete(id int) error
	FindByCode(code string) ([]Grade, error)
}
