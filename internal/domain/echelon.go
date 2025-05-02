package domain

import (
	"time"
)

type Echelon struct {
	ID         int       `json:"id" db:"id"`
	Code       string    `json:"code" db:"code"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}
type EchelonInterface interface {
	FindAll() ([]Echelon, error)
	FindByID(id int) (*Echelon, error)
	Save(echelon *Echelon) (*Echelon, error)
	Update(echelon *Echelon) (*Echelon, error)
	Delete(id int) error
	FindByCode(code string) (*Echelon, error)
}
