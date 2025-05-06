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

/*
-- DDL

create table achmadnr.employee_assignments(

	employee_id uuid not null,
	unit_id int not null,
	position_id int not null,
	is_active boolean not null,
	assigned_at timestamp default now(),
	created_at timestamp default now(),
	modified_at timestamp default now(),
	primary key(employee_id, position_id, unit_id),
	foreign key(employee_id) references achmadnr.employees(id),
	foreign key(unit_id) references achmadnr.units(id),
	foreign key(position_id) references achmadnr.positions(id)

);

	type EmployeeAssignment struct {
		EmployeeID string    `json:"employee_id" db:"employee_id"`
		UnitID     int       `json:"unit_id" db:"unit_id"`
		PositionID int       `json:"position_id" db:"position_id"`
		IsActive   bool      `json:"is_active" db:"is_active"`
		AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	}

	type EmployeeAssignmentResponse struct {
		EmployeeID   string    `json:"employee_id"`
		UnitID       int       `json:"unit_id"`
		PositionID   int       `json:"position_id"`
		EmployeeName string    `json:"employee_name"`
		UnitName     string    `json:"unit_name"`
		PositionName string    `json:"position_name"`
		IsActive     bool      `json:"is_active"`
		AssignedAt   time.Time `json:"assigned_at"`
	}

	type EmployeeAssignmentInterface interface {
		Save(employeeAssignment *EmployeeAssignment) (*EmployeeAssignment, error)
		Update(employeeAssignment *EmployeeAssignment) (*EmployeeAssignment, error)
		Delete(employeeID string, unitID int, positionID int) error
		FindAll() ([]EmployeeAssignmentResponse, error)
		FindByID(employeeID string, unitID int, positionID int) (*EmployeeAssignmentResponse, error)
		FindByEmployeeID(employeeID string) (*EmployeeAssignmentResponse, error)
		FindByUnitID(unitID int) ([]EmployeeAssignmentResponse, error)
	}
*/

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

func (r *EmployeeAssignmentRepository) FindAll() ([]domain.EmployeeAssignmentResponse, error) {
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.name as employee_name, u.name as unit_name, p.name as position_name
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
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, ea.created_at, ea.modified_at, e.name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.employee_id = $1 AND ea.unit_id = $2 AND ea.position_id = $3`
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
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.employee_id = $1`
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
	query := `SELECT ea.employee_id, ea.unit_id, ea.position_id, ea.is_active, ea.assigned_at, e.name as employee_name, u.name as unit_name, p.name as position_name
	FROM achmadnr.employee_assignments ea
	INNER JOIN achmadnr.employees e ON ea.employee_id = e.id
	INNER JOIN achmadnr.units u ON ea.unit_id = u.id
	INNER JOIN achmadnr.positions p ON ea.position_id = p.id
	WHERE ea.unit_id = $1`
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
