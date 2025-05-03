package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type UnitRepository struct {
	db *sql.DB
}

func NewUnitRepository(db *sql.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

func (r *UnitRepository) FindAll() ([]domain.Unit, error) {
	query := `SELECT id, name, address, description, created_at, modified_at FROM achmadnr.units`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []domain.Unit
	for rows.Next() {
		var unit domain.Unit
		if err := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.Address,
			&unit.Description,
			&unit.CreatedAt,
			&unit.ModifiedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}
func (r *UnitRepository) FindByID(id int) (*domain.Unit, error) {
	query := `SELECT id, name, address, description, created_at, modified_at FROM achmadnr.units WHERE id = $1`
	unit := &domain.Unit{}
	err := r.db.QueryRow(query, id).Scan(
		&unit.ID,
		&unit.Name,
		&unit.Address,
		&unit.Description,
		&unit.CreatedAt,
		&unit.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
func (r *UnitRepository) Save(unit *domain.Unit) (*domain.Unit, error) {
	query := `INSERT INTO achmadnr.units (name, address, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, unit.Name, unit.Address, unit.Description).Scan(&unit.ID)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
func (r *UnitRepository) Update(unit *domain.Unit) (*domain.Unit, error) {
	query := `UPDATE achmadnr.units SET name = $1, address = $2, description = $3, modified_at = now() WHERE id = $4`
	_, err := r.db.Exec(query, unit.Name, unit.Address, unit.Description, unit.ID)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
func (r *UnitRepository) Delete(id int) error {
	query := `DELETE FROM achmadnr.units WHERE id = $1`
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
func (r *UnitRepository) FindByName(name string) ([]domain.Unit, error) {
	query := `SELECT id, name, address, description, created_at, modified_at FROM achmadnr.units WHERE name ILIKE $1`
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []domain.Unit
	for rows.Next() {
		var unit domain.Unit
		if err := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.Address,
			&unit.Description,
			&unit.CreatedAt,
			&unit.ModifiedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}
