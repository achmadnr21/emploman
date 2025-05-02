package domain

import (
	"time"
)

/*
create table achmadnr.positions(
	id SERIAL primary key,
	name varchar(255) not null,
	created_at timestamp default now(),
	modified_at timestamp default now()
);

*/

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
	FindByName(name string) (*Position, error)
}
