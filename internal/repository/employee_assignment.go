package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type EmployeeAssignmentRepository struct {
	db *sql.DB
}

func NewEmployeeAssignmentRepository(db *sql.DB) *EmployeeAssignmentRepository {
	return &EmployeeAssignmentRepository{
		db: db,
	}
}

func (r *EmployeeAssignmentRepository) FindAll() ([]domain.EmployeeAssignment, error) {
	query := `SELECT employee_id, unit_id, position_id, is_active, assigned_at, created_at, modified_at FROM achmadnr.employee_assignments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeAssignments []domain.EmployeeAssignment
	for rows.Next() {
		var employeeAssignment domain.EmployeeAssignment
		if err := rows.Scan(&employeeAssignment.EmployeeID, &employeeAssignment.UnitID, &employeeAssignment.PositionID, &employeeAssignment.IsActive, &employeeAssignment.AssignedAt, &employeeAssignment.CreatedAt, &employeeAssignment.ModifiedAt); err != nil {
			return nil, err
		}
		employeeAssignments = append(employeeAssignments, employeeAssignment)
	}

	return employeeAssignments, nil
}
func (r *EmployeeAssignmentRepository) FindByID(employeeID string, unitID int, positionID int) (*domain.EmployeeAssignment, error) {
	query := `SELECT employee_id, unit_id, position_id, is_active, assigned_at, created_at, modified_at FROM achmadnr.employee_assignments WHERE employee_id = $1 AND unit_id = $2 AND position_id = $3`
	row := r.db.QueryRow(query, employeeID, unitID, positionID)

	var employeeAssignment *domain.EmployeeAssignment = &domain.EmployeeAssignment{}
	if err := row.Scan(&employeeAssignment.EmployeeID, &employeeAssignment.UnitID, &employeeAssignment.PositionID, &employeeAssignment.IsActive, &employeeAssignment.AssignedAt, &employeeAssignment.CreatedAt, &employeeAssignment.ModifiedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return employeeAssignment, nil
}
func (r *EmployeeAssignmentRepository) Save(employeeAssignment *domain.EmployeeAssignment) (*domain.EmployeeAssignment, error) {
	query := `INSERT INTO achmadnr.employee_assignments (employee_id, unit_id, position_id, is_active) VALUES ($1, $2, $3, $4) RETURNING employee_id, assigned_at, created_at, modified_at`
	err := r.db.QueryRow(query,
		employeeAssignment.EmployeeID,
		employeeAssignment.UnitID,
		employeeAssignment.PositionID,
		employeeAssignment.IsActive,
	).Scan(&employeeAssignment.EmployeeID, &employeeAssignment.AssignedAt, &employeeAssignment.CreatedAt, &employeeAssignment.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return employeeAssignment, nil
}
func (r *EmployeeAssignmentRepository) Update(employeeAssignment *domain.EmployeeAssignment) (*domain.EmployeeAssignment, error) {
	query := `UPDATE achmadnr.employee_assignments SET unit_id = $1, position_id = $2, is_active = $3, assigned_at = $4, modified_at = NOW() WHERE employee_id = $5 AND unit_id = $6 AND position_id = $7`
	_, err := r.db.Exec(query,
		employeeAssignment.UnitID,
		employeeAssignment.PositionID,
		employeeAssignment.IsActive,
		employeeAssignment.AssignedAt,
		employeeAssignment.EmployeeID,
		employeeAssignment.UnitID,
		employeeAssignment.PositionID,
	)
	if err != nil {
		return nil, err
	}

	return employeeAssignment, nil
}
func (r *EmployeeAssignmentRepository) Delete(employeeID string, unitID int, positionID int) error {
	query := `DELETE FROM achmadnr.employee_assignments WHERE employee_id = $1 AND unit_id = $2 AND position_id = $3`
	res, err := r.db.Exec(query, employeeID, unitID, positionID)
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
func (r *EmployeeAssignmentRepository) FindByEmployeeID(employeeID string) (*domain.EmployeeAssignment, error) {
	query := `SELECT employee_id, unit_id, position_id, is_active, assigned_at, created_at, modified_at FROM achmadnr.employee_assignments WHERE employee_id = $1`
	row := r.db.QueryRow(query, employeeID)

	var employeeAssignment *domain.EmployeeAssignment = &domain.EmployeeAssignment{}
	if err := row.Scan(&employeeAssignment.EmployeeID, &employeeAssignment.UnitID, &employeeAssignment.PositionID, &employeeAssignment.IsActive, &employeeAssignment.AssignedAt, &employeeAssignment.CreatedAt, &employeeAssignment.ModifiedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return employeeAssignment, nil
}
func (r *EmployeeAssignmentRepository) FindByUnitID(unitID int) ([]domain.EmployeeAssignment, error) {
	query := `SELECT employee_id, unit_id, position_id, is_active, assigned_at, created_at, modified_at FROM achmadnr.employee_assignments WHERE unit_id = $1`
	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeAssignments []domain.EmployeeAssignment
	for rows.Next() {
		var employeeAssignment domain.EmployeeAssignment
		if err := rows.Scan(&employeeAssignment.EmployeeID, &employeeAssignment.UnitID, &employeeAssignment.PositionID, &employeeAssignment.IsActive, &employeeAssignment.AssignedAt, &employeeAssignment.CreatedAt, &employeeAssignment.ModifiedAt); err != nil {
			return nil, err
		}
		employeeAssignments = append(employeeAssignments, employeeAssignment)
	}

	return employeeAssignments, nil
}
