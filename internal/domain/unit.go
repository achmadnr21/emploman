package domain

import (
	"time"
)

/*
create table achmadnr.units(
	id SERIAL primary key,
	name varchar(255) not null,
	address text not null,
	description text default 'no desc',
	created_at timestamp default now(),
	modified_at timestamp default now()
);
*/

type Unit struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ModifiedAt  time.Time `json:"modified_at" db:"modified_at"`
}

type UnitInterface interface {
	FindAll() ([]Unit, error)
	FindByID(id int) (*Unit, error)
	Save(unit *Unit) (*Unit, error)
	Update(unit *Unit) (*Unit, error)
	Delete(id int) error
	FindByName(name string) (*Unit, error)
}
