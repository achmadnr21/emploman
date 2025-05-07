package domain

import (
	"time"
)

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
	TransactionalAssignment(employeeAssignment *EmployeeAssignment) error
	Deactivate(employeeID string, unitID int, positionID int) error
	FindAll() ([]EmployeeAssignmentResponse, error)
	FindByID(employeeID string, unitID int, positionID int) (*EmployeeAssignmentResponse, error)
	FindByEmployeeID(employeeID string) (*EmployeeAssignmentResponse, error)
	FindByUnitID(unitID int) ([]EmployeeAssignmentResponse, error)
}
