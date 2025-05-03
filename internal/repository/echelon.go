package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type EchelonRepository struct {
	db *sql.DB
}

func NewEchelonRepository(db *sql.DB) *EchelonRepository {
	return &EchelonRepository{
		db: db,
	}
}

func (r *EchelonRepository) FindAll() ([]domain.Echelon, error) {
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
func (r *EchelonRepository) FindByID(id int) (*domain.Echelon, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.echelons WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var echelon domain.Echelon
	if err := row.Scan(&echelon.ID, &echelon.Code, &echelon.CreatedAt, &echelon.ModifiedAt); err != nil {
		return nil, err
	}
	return &echelon, nil
}
func (r *EchelonRepository) Save(echelon *domain.Echelon) (*domain.Echelon, error) {
	query := `INSERT INTO achmadnr.echelons (id, code) VALUES ($1, $2)`
	_, err := r.db.Exec(query, echelon.ID, echelon.Code)
	if err != nil {
		return nil, err
	}
	return echelon, nil
}
func (r *EchelonRepository) Update(echelon *domain.Echelon) (*domain.Echelon, error) {
	query := `UPDATE achmadnr.echelons SET code = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, echelon.Code, echelon.ID)
	if err != nil {
		return nil, err
	}
	return echelon, nil
}
func (r *EchelonRepository) Delete(id int) error {
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
func (r *EchelonRepository) FindByCode(code string) ([]domain.Echelon, error) {
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
