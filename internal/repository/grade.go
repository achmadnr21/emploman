package repository

import (
	"database/sql"

	"github.com/achmadnr21/emploman/internal/domain"
)

type GradeRepository struct {
	db *sql.DB
}

func NewGradeRepository(db *sql.DB) *GradeRepository {
	return &GradeRepository{
		db: db,
	}
}

func (r *GradeRepository) FindAll() ([]domain.Grade, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.grades`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var grades []domain.Grade
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Code, &grade.CreatedAt, &grade.ModifiedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return grades, nil
}
func (r *GradeRepository) FindByID(id int) (*domain.Grade, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.grades WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var grade domain.Grade
	if err := row.Scan(&grade.ID, &grade.Code, &grade.CreatedAt, &grade.ModifiedAt); err != nil {
		return nil, err
	}
	return &grade, nil
}
func (r *GradeRepository) Save(grade *domain.Grade) (*domain.Grade, error) {
	query := `INSERT INTO achmadnr.grades (id, code) VALUES ($1, $2)`
	_, err := r.db.Exec(query, grade.ID, grade.Code)
	if err != nil {
		return nil, err
	}
	return grade, nil
}
func (r *GradeRepository) Update(grade *domain.Grade) (*domain.Grade, error) {
	query := `UPDATE achmadnr.grades SET code = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, grade.Code, grade.ID)
	if err != nil {
		return nil, err
	}
	return grade, nil
}
func (r *GradeRepository) Delete(id int) error {
	query := `DELETE FROM achmadnr.grades WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *GradeRepository) FindByCode(code string) (*domain.Grade, error) {
	query := `SELECT id, code, created_at, modified_at FROM achmadnr.grades WHERE code = $1`
	row := r.db.QueryRow(query, code)
	var grade domain.Grade
	if err := row.Scan(&grade.ID, &grade.Code, &grade.CreatedAt, &grade.ModifiedAt); err != nil {
		return nil, err
	}
	return &grade, nil
}
