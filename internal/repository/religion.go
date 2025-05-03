package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type ReligionRepository struct {
	db *sql.DB
}

func NewReligionRepository(db *sql.DB) *ReligionRepository {
	return &ReligionRepository{
		db: db,
	}
}

func (r *ReligionRepository) FindAll() ([]domain.Religion, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.religions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var religions []domain.Religion
	for rows.Next() {
		var religion domain.Religion
		if err := rows.Scan(&religion.ID, &religion.Name, &religion.CreatedAt, &religion.ModifiedAt); err != nil {
			return nil, err
		}
		religions = append(religions, religion)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return religions, nil
}

func (r *ReligionRepository) FindByID(id string) (*domain.Religion, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.religions WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var religion domain.Religion
	if err := row.Scan(&religion.ID, &religion.Name, &religion.CreatedAt, &religion.ModifiedAt); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}
	return &religion, nil
}
func (r *ReligionRepository) Save(religion *domain.Religion) (*domain.Religion, error) {
	query := `INSERT INTO achmadnr.religions (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, religion.ID, religion.Name)
	if err != nil {
		return nil, err
	}
	return religion, nil
}
func (r *ReligionRepository) Update(religion *domain.Religion) (*domain.Religion, error) {
	query := `UPDATE achmadnr.religions SET name = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, religion.Name, religion.ID)
	if err != nil {
		return nil, err
	}
	return religion, nil
}
func (r *ReligionRepository) Delete(id string) error {
	query := `DELETE FROM achmadnr.religions WHERE id = $1`
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
func (r *ReligionRepository) FindByName(name string) ([]domain.Religion, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.religions WHERE name ILIKE $1`
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var religions []domain.Religion
	for rows.Next() {
		var religion domain.Religion
		if err := rows.Scan(&religion.ID, &religion.Name, &religion.CreatedAt, &religion.ModifiedAt); err != nil {
			return nil, err
		}
		religions = append(religions, religion)
	}
	return religions, nil
}
