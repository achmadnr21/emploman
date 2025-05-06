package domain

import (
	"time"
)

type Role struct {
	ID                        string    `json:"id" db:"id"`
	Name                      string    `json:"name" db:"name"`
	Level                     int       `json:"level" db:"level"`
	Description               string    `json:"description" db:"description"`
	CanAddRole                bool      `json:"can_add_role" db:"can_add_role"`
	CanAddEmployee            bool      `json:"can_add_employee" db:"can_add_employee"`
	CanAddUnit                bool      `json:"can_add_unit" db:"can_add_unit"`
	CanAddPosition            bool      `json:"can_add_position" db:"can_add_position"`
	CanAddEchelon             bool      `json:"can_add_echelon" db:"can_add_echelon"`
	CanAddReligion            bool      `json:"can_add_religion" db:"can_add_religion"`
	CanAddGrade               bool      `json:"can_add_grade" db:"can_add_grade"`
	CanAssignEmployeeInternal bool      `json:"can_assign_employee_internal" db:"can_assign_employee_internal"`
	CanAssignEmployeeGlobal   bool      `json:"can_assign_employee_global" db:"can_assign_employee_global"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	ModifiedAt                time.Time `json:"modified_at" db:"modified_at"`
}

/*
CREATE TABLE achmadnr.role_promotions (

	promoter_role_id CHAR(3) not null,
	from_role_id CHAR(3) NOT NULL,
	to_role_id CHAR(3) NOT NULL,
	PRIMARY KEY (promoter_role_id, from_role_id, to_role_id),
	foreign key (promoter_role_id) references achmadnr.roles(id) on delete cascade,
	FOREIGN KEY (from_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE,
	FOREIGN KEY (to_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE

);
*/
type RolePromotion struct {
	PromoterRoleID string `json:"promoter_role_id" db:"promoter_role_id"`
	FromRoleID     string `json:"from_role_id" db:"from_role_id"`
	ToRoleID       string `json:"to_role_id" db:"to_role_id"`
}

type RoleInterface interface {
	FindByUserID(id string) (*Role, error)
	FindAll() ([]Role, error)
	FindByID(id string) (*Role, error)
	Save(role *Role) (*Role, error)
	Update(role *Role) (*Role, error)
	Delete(id string) error
	FindByName(name string) (*Role, error)
	FindPromoteRole(promoterRoleID string) ([]RolePromotion, error)
}
