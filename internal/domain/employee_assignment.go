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

type EmployeeAssignmentInterface interface {
	FindAll() ([]EmployeeAssignment, error)
	FindByID(employeeID string, unitID int, positionID int) (*EmployeeAssignment, error)
	Save(employeeAssignment *EmployeeAssignment) (*EmployeeAssignment, error)
	Update(employeeAssignment *EmployeeAssignment) (*EmployeeAssignment, error)
	Delete(employeeID string, unitID int, positionID int) error
	FindByEmployeeID(employeeID string) (*EmployeeAssignment, error)
	FindByUnitID(unitID int) ([]EmployeeAssignment, error)
}
