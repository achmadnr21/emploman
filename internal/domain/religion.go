package domain

import (
	"time"
)

type Religion struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

type ReligionInterface interface {
	Save(religion *Religion) (*Religion, error)
	Update(religion *Religion) (*Religion, error)
	Delete(id string) error
	FindAll() ([]Religion, error)
	FindByID(id string) (*Religion, error)
	FindByName(name string) ([]Religion, error)
}
