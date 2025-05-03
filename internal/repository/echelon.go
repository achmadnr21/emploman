package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type EcholonRepository struct {
	db *sql.DB
}

func NewEcholonRepository(db *sql.DB) *EcholonRepository {
	return &EcholonRepository{
		db: db,
	}
}

/*
-- DDL

create table achmadnr.echelons(
	id SERIAL unique primary key,
	code varchar(100),
	created_at timestamp default now(),
	modified_at timestamp default now()
);

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
	FindByCode(code string) ([]Echelon, error)
}

*/

func (r *EcholonRepository) FindAll() ([]domain.Echelon, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.echelons`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var echelons []domain.Echelon
	for rows.Next() {
		var echelon domain.Echelon
		err := rows.Scan(&echelon.ID, &echelon.Code, &echelon.CreatedAt, &echelon.ModifiedAt)
		if err != nil {
			return nil, err
		}
		echelons = append(echelons, echelon)
	}
	return echelons, nil
}
func (r *EcholonRepository) FindByID(id int) (*domain.Echelon, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.echelons WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var echelon domain.Echelon
	if err := row.Scan(&echelon.ID, &echelon.Code, &echelon.CreatedAt, &echelon.ModifiedAt); err != nil {
		return nil, err
	}
	return &echelon, nil
}
func (r *EcholonRepository) Save(echelon *domain.Echelon) (*domain.Echelon, error) {
	query := `INSERT INTO achmadnr.echelons (id, code) VALUES ($1, $2)`
	_, err := r.db.Exec(query, echelon.ID, echelon.Code)
	if err != nil {
		return nil, err
	}
	return echelon, nil
}
func (r *EcholonRepository) Update(echelon *domain.Echelon) (*domain.Echelon, error) {
	query := `UPDATE achmadnr.echelons SET code = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, echelon.Code, echelon.ID)
	if err != nil {
		return nil, err
	}
	return echelon, nil
}
func (r *EcholonRepository) Delete(id int) error {
	query := `DELETE FROM achmadnr.echelons WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}
func (r *EcholonRepository) FindByCode(code string) ([]domain.Echelon, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.echelons WHERE code ILIKE $1`
	rows, err := r.db.Query(query, "%"+code+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var echelons []domain.Echelon
	for rows.Next() {
		var echelon domain.Echelon
		err := rows.Scan(&echelon.ID, &echelon.Code, &echelon.CreatedAt, &echelon.ModifiedAt)
		if err != nil {
			return nil, err
		}
		echelons = append(echelons, echelon)
	}
	return echelons, nil
}
