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

func (r *EmployeeAssignmentRepository) TransactionalAssignment(employeeAssignment *domain.EmployeeAssignment) error {

	tx, err := r.db.Begin() // e.db = *sql.DB, pastikan sudah ada
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 1. Nonaktifkan assignment aktif (jika ada)
	_, err = tx.Exec(`
		UPDATE achmadnr.employee_assignments
		SET is_active = FALSE, assigned_at = NOW()
		WHERE employee_id = $1 AND is_active = TRUE
	`, employeeAssignment.EmployeeID)
	if err != nil {
		return err
	}

	// 2. Coba aktifkan kembali assignment lama
	var dummy int
	err = tx.QueryRow(`
		WITH updated AS (
			UPDATE achmadnr.employee_assignments
			SET is_active = TRUE, assigned_at = NOW()
			WHERE employee_id = $1 AND unit_id = $2 AND position_id = $3
			RETURNING 1
		)
		SELECT 1 FROM updated LIMIT 1
	`, employeeAssignment.EmployeeID, employeeAssignment.UnitID, employeeAssignment.PositionID).Scan(&dummy)

	if err == sql.ErrNoRows {
		// 3. Tidak ditemukan, insert baru
		_, err = tx.Exec(`
			INSERT INTO achmadnr.employee_assignments (
				employee_id, unit_id, position_id, is_active, assigned_at
			) VALUES ($1, $2, $3, TRUE, NOW())
		`, employeeAssignment.EmployeeID, employeeAssignment.UnitID, employeeAssignment.PositionID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil

}

func (r *EmployeeAssignmentRepository) Deactivate(employeeID string, unitID int, positionID int) error {
	// update is_active to false
	query := `UPDATE achmadnr.employee_assignments SET is_active = FALSE, assigned_at = NOW() WHERE employee_id = $1 AND unit_id = $2 AND position_id = $3`
	_, err := r.db.Exec(query, employeeID, unitID, positionID)
	if err != nil {
		return fmt.Errorf("failed to deactivate employee assignment: %w", err)
	}

	return nil
}

func (r *EmployeeAssignmentRepository) FindAll() ([]domain.EmployeeAssignmentResponse, error) {
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.full_name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeAssignments []domain.EmployeeAssignmentResponse
	for rows.Next() {
		var employeeAssignment domain.EmployeeAssignmentResponse
		if err := rows.Scan(
			&employeeAssignment.EmployeeID,
			&employeeAssignment.UnitID,
			&employeeAssignment.PositionID,
			&employeeAssignment.IsActive,
			&employeeAssignment.AssignedAt,
			&employeeAssignment.EmployeeName,
			&employeeAssignment.UnitName,
			&employeeAssignment.PositionName); err != nil {
			return nil, err
		}
		employeeAssignments = append(employeeAssignments, employeeAssignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return employeeAssignments, nil
}
func (r *EmployeeAssignmentRepository) FindByID(employeeID string, unitID int, positionID int) (*domain.EmployeeAssignmentResponse, error) {
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, ea.created_at, ea.modified_at, e.full_name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.employee_id = $1 AND ea.unit_id = $2 AND ea.position_id = $3
	ORDER BY ea.assigned_at DESC
	LIMIT 1`
	row := r.db.QueryRow(query, employeeID, unitID, positionID)
	var employeeAssignment domain.EmployeeAssignmentResponse
	if err := row.Scan(
		&employeeAssignment.EmployeeID,
		&employeeAssignment.UnitID,
		&employeeAssignment.PositionID,
		&employeeAssignment.IsActive,
		&employeeAssignment.AssignedAt,
		&employeeAssignment.EmployeeName,
		&employeeAssignment.UnitName,
		&employeeAssignment.PositionName); err != nil {
		return nil, err
	}
	return &employeeAssignment, nil
}

func (r *EmployeeAssignmentRepository) FindByEmployeeID(employeeID string) (*domain.EmployeeAssignmentResponse, error) {
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.full_name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.employee_id = $1 
	ORDER BY ea.assigned_at DESC
	LIMIT 1`
	row := r.db.QueryRow(query, employeeID)
	var employeeAssignment domain.EmployeeAssignmentResponse
	if err := row.Scan(
		&employeeAssignment.EmployeeID,
		&employeeAssignment.UnitID,
		&employeeAssignment.PositionID,
		&employeeAssignment.IsActive,
		&employeeAssignment.AssignedAt,
		&employeeAssignment.EmployeeName,
		&employeeAssignment.UnitName,
		&employeeAssignment.PositionName); err != nil {
		return nil, err
	}
	return &employeeAssignment, nil
}
func (r *EmployeeAssignmentRepository) FindByUnitID(unitID int) ([]domain.EmployeeAssignmentResponse, error) {
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.full_name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.unit_id = $1
	ORDER BY ea.assigned_at DESC`
	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employeeAssignments []domain.EmployeeAssignmentResponse
	for rows.Next() {
		var employeeAssignment domain.EmployeeAssignmentResponse
		if err := rows.Scan(
			&employeeAssignment.EmployeeID,
			&employeeAssignment.UnitID,
			&employeeAssignment.PositionID,
			&employeeAssignment.IsActive,
			&employeeAssignment.AssignedAt,
			&employeeAssignment.EmployeeName,
			&employeeAssignment.UnitName,
			&employeeAssignment.PositionName); err != nil {
			return nil, err
		}
		employeeAssignments = append(employeeAssignments, employeeAssignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return employeeAssignments, nil
}
