package repository

import (
	"database/sql"

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
func (r *RoleRepository) FindAll() ([]domain.Role, error) {
	query := `SELECT 
	id, name, description, can_add_role, can_add_user, can_add_unit, can_add_position, 
	can_add_echelon, can_add_religion, can_add_grade, can_assign_employee, can_assign_employee_global, 
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
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CanAddRole, &role.CanAddUser, &role.CanAddUnit, &role.CanAddPosition, &role.CanAddEchelon, &role.CanAddReligion, &role.CanAddGrade, &role.CanAssignEmployee, &role.CanAssignEmployeeGlobal, &role.CreatedAt, &role.ModifiedAt); err != nil {
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
	query := `SELECT id, name, description, can_add_role, can_add_user, can_add_unit, 
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee, 
	can_assign_employee_global, created_at, modified_at FROM achmadnr.roles WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var role domain.Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CanAddRole, &role.CanAddUser, &role.CanAddUnit, &role.CanAddPosition, &role.CanAddEchelon, &role.CanAddReligion, &role.CanAddGrade, &role.CanAssignEmployee, &role.CanAssignEmployeeGlobal, &role.CreatedAt, &role.ModifiedAt); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}
	return &role, nil
}
func (r *RoleRepository) Save(role *domain.Role) (*domain.Role, error) {
	query := `INSERT INTO achmadnr.roles (id, name, description, can_add_role, can_add_user, can_add_unit,
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee,
	can_assign_employee_global)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(query,
		role.ID,
		role.Name,
		role.Description,
		role.CanAddRole,
		role.CanAddUser,
		role.CanAddUnit,
		role.CanAddPosition,
		role.CanAddEchelon,
		role.CanAddReligion,
		role.CanAddGrade,
		role.CanAssignEmployee,
		role.CanAssignEmployeeGlobal)
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (r *RoleRepository) Update(role *domain.Role) (*domain.Role, error) {
	query := `UPDATE achmadnr.roles SET name = $1, description = $2, can_add_role = $3, can_add_user = $4,
	can_add_unit = $5, can_add_position = $6, can_add_echelon = $7, can_add_religion = $8,
	can_add_grade = $9, can_assign_employee = $10, can_assign_employee_global = $11,
	modified_at = now() WHERE id = $12`
	_, err := r.db.Exec(query,
		role.Name,
		role.Description,
		role.CanAddRole,
		role.CanAddUser,
		role.CanAddUnit,
		role.CanAddPosition,
		role.CanAddEchelon,
		role.CanAddReligion,
		role.CanAddGrade,
		role.CanAssignEmployee,
		role.CanAssignEmployeeGlobal,
		role.ID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Delete(id string) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *RoleRepository) FindByName(name string) (*domain.Role, error) {
	query := `SELECT id, name, description, can_add_role, can_add_user, can_add_unit, 
	can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee, 
	can_assign_employee_global, created_at, modified_at FROM achmadnr.roles WHERE lower(name) LIKE lower($1)`
	row := r.db.QueryRow(query, "%"+name+"%")
	var role domain.Role
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CanAddRole,
		&role.CanAddUser,
		&role.CanAddUnit,
		&role.CanAddPosition,
		&role.CanAddEchelon,
		&role.CanAddReligion,
		&role.CanAddGrade,
		&role.CanAssignEmployee,
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
