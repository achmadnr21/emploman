package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

/*
	type RoleInterface interface {
		FindAll() ([]Role, error)
		FindByID(id string) (*Role, error)
		Save(role *Role) (*Role, error)
		Update(role *Role) (*Role, error)
		Delete(id string) error
		FindByName(name string) (*Role, error)
		FindPromoteRole(promoterRoleID, fromRoleID, toRoleID string) (*RolePromotion, error)
	}
*/
func (r *RoleRepository) FindAll() ([]domain.Role, error) {
	query := `SELECT 
	id, name, level, description, can_add_role, can_add_employee, can_add_unit, can_add_position, 
	can_add_echelon, can_add_religion, can_add_grade, can_assign_employee_internal, can_assign_employee_global, 
	created_at, modified_at 
	FROM achmadnr.roles`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Level, &role.Description, &role.CanAddRole, &role.CanAddEmployee, &role.CanAddUnit, &role.CanAddPosition, &role.CanAddEchelon, &role.CanAddReligion, &role.CanAddGrade, &role.CanAssignEmployeeInternal, &role.CanAssignEmployeeGlobal, &role.CreatedAt, &role.ModifiedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
func (r *RoleRepository) FindByID(id string) (*domain.Role, error) {
	query := `SELECT id, name, level, description, can_add_role, can_add_employee, can_add_unit, 
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee_internal, 
	can_assign_employee_global, created_at, modified_at FROM achmadnr.roles WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var role domain.Role
	if err := row.Scan(&role.ID, &role.Name, &role.Level, &role.Description, &role.CanAddRole, &role.CanAddEmployee, &role.CanAddUnit, &role.CanAddPosition, &role.CanAddEchelon, &role.CanAddReligion, &role.CanAddGrade, &role.CanAssignEmployeeInternal, &role.CanAssignEmployeeGlobal, &role.CreatedAt, &role.ModifiedAt); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}
	return &role, nil
}
func (r *RoleRepository) Save(role *domain.Role) (*domain.Role, error) {
	query := `INSERT INTO achmadnr.roles (id, name, level, description, can_add_role, can_add_employee, can_add_unit,
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee_internal,
	can_assign_employee_global)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(query,
		role.ID,
		role.Name,
		role.Level,
		role.Description,
		role.CanAddRole,
		role.CanAddEmployee,
		role.CanAddUnit,
		role.CanAddPosition,
		role.CanAddEchelon,
		role.CanAddReligion,
		role.CanAddGrade,
		role.CanAssignEmployeeInternal,
		role.CanAssignEmployeeGlobal)
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (r *RoleRepository) Update(role *domain.Role) (*domain.Role, error) {
	query := `UPDATE achmadnr.roles SET name = $1, level = $2, description = $3, can_add_role = $4, can_add_employee = $5,
	can_add_unit = $6, can_add_position = $7, can_add_echelon = $8, can_add_religion = $9,
	can_add_grade = $10, can_assign_employee_internal = $11, can_assign_employee_global = $12,
	modified_at = now() WHERE id = $13`
	_, err := r.db.Exec(query,
		role.Name,
		role.Level,
		role.Description,
		role.CanAddRole,
		role.CanAddEmployee,
		role.CanAddUnit,
		role.CanAddPosition,
		role.CanAddEchelon,
		role.CanAddReligion,
		role.CanAddGrade,
		role.CanAssignEmployeeInternal,
		role.CanAssignEmployeeGlobal,
		role.ID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Delete(id string) error {
	query := `DELETE FROM roles WHERE id = $1`
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
func (r *RoleRepository) FindByName(name string) (*domain.Role, error) {
	query := `SELECT id, name, level, description, can_add_role, can_add_employee, can_add_unit, 
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee_internal, 
	can_assign_employee_global, created_at, modified_at FROM achmadnr.roles WHERE name ILIKE $1`
	row := r.db.QueryRow(query, "%"+name+"%")
	var role domain.Role
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
		&role.CanAddRole,
		&role.CanAddEmployee,
		&role.CanAddUnit,
		&role.CanAddPosition,
		&role.CanAddEchelon,
		&role.CanAddReligion,
		&role.CanAddGrade,
		&role.CanAssignEmployeeInternal,
		&role.CanAssignEmployeeGlobal,
		&role.CreatedAt,
		&role.ModifiedAt)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindPromoteRole(promoterRoleID, fromRoleID, toRoleID string) (*domain.RolePromotion, error) {
	query := `SELECT promoter_role_id, from_role_id, to_role_id FROM achmadnr.role_promotions WHERE promoter_role_id = $1 AND from_role_id = $2 AND to_role_id = $3`
	row := r.db.QueryRow(query, promoterRoleID, fromRoleID, toRoleID)
	var rolePromotion domain.RolePromotion
	err := row.Scan(&rolePromotion.PromoterRoleID, &rolePromotion.FromRoleID, &rolePromotion.ToRoleID)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}
	return &rolePromotion, nil
}
