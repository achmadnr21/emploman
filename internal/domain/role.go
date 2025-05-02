package domain

import (
	"time"
)

type Role struct {
	ID                      string    `json:"id" db:"id"`
	Name                    string    `json:"name" db:"name"`
	Description             string    `json:"description" db:"description"`
	CanAddRole              bool      `json:"can_add_role" db:"can_add_role"`
	CanAddUser              bool      `json:"can_add_user" db:"can_add_user"`
	CanAddUnit              bool      `json:"can_add_unit" db:"can_add_unit"`
	CanAddPosition          bool      `json:"can_add_position" db:"can_add_position"`
	CanAddEchelon           bool      `json:"can_add_echelon" db:"can_add_echelon"`
	CanAddReligion          bool      `json:"can_add_religion" db:"can_add_religion"`
	CanAddGrade             bool      `json:"can_add_grade" db:"can_add_grade"`
	CanAssignEmployee       bool      `json:"can_assign_employee" db:"can_assign_employee"`
	CanAssignEmployeeGlobal bool      `json:"can_assign_employee_global" db:"can_assign_employee_global"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	ModifiedAt              time.Time `json:"modified_at" db:"modified_at"`
}
type RoleInterface interface {
	FindAll() ([]Role, error)
	FindByID(id string) (*Role, error)
	Save(role *Role) (*Role, error)
	Update(role *Role) (*Role, error)
	Delete(id string) error
}
