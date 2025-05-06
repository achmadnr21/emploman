package domain

import (
	"time"
)

type Position struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}
type PositionInterface interface {
	FindAll() ([]Position, error)
	FindByID(id int) (*Position, error)
	Save(position *Position) (*Position, error)
	Update(position *Position) (*Position, error)
	Delete(id int) error
	FindByName(name string) ([]Position, error)
	Search(query string) ([]Position, error)
}
