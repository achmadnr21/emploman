package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type PositionRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) *PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

func (r *PositionRepository) FindAll() ([]domain.Position, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.positions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var position domain.Position
		if err := rows.Scan(
			&position.ID,
			&position.Name,
			&position.CreatedAt,
			&position.ModifiedAt,
		); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	return positions, nil
}
func (r *PositionRepository) FindByID(id int) (*domain.Position, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.positions WHERE id = $1`
	position := &domain.Position{}
	err := r.db.QueryRow(query, id).Scan(
		&position.ID,
		&position.Name,
		&position.CreatedAt,
		&position.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return position, nil
}
func (r *PositionRepository) Save(position *domain.Position) (*domain.Position, error) {
	query := `INSERT INTO achmadnr.positions (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(query, position.Name).Scan(&position.ID)
	if err != nil {
		return nil, err
	}
	return position, nil
}
func (r *PositionRepository) Update(position *domain.Position) (*domain.Position, error) {
	query := `UPDATE achmadnr.positions SET name = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, position.Name, position.ID)
	if err != nil {
		return nil, err
	}
	return position, nil
}
func (r *PositionRepository) Delete(id int) error {
	query := `DELETE FROM achmadnr.positions WHERE id = $1`
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
func (r *PositionRepository) FindByName(name string) ([]domain.Position, error) {
	query := `SELECT id, name, created_at, modified_at FROM achmadnr.positions WHERE name ILIKE $1`
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var positions []domain.Position
	for rows.Next() {
		var position domain.Position
		if err := rows.Scan(
			&position.ID,
			&position.Name,
			&position.CreatedAt,
			&position.ModifiedAt,
		); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

func (r *PositionRepository) Search(query string) ([]domain.Position, error) {
	query = fmt.Sprintf("%%%s%%", query)
	sqlQuery := `SELECT id, name, created_at, modified_at FROM achmadnr.positions WHERE name ILIKE $1`
	rows, err := r.db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var position domain.Position
		if err := rows.Scan(
			&position.ID,
			&position.Name,
			&position.CreatedAt,
			&position.ModifiedAt,
		); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	return positions, nil
}
