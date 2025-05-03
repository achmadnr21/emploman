package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) FindAll() ([]domain.Employee, error) {
	query := `SELECT id, role_id, nip, password, full_name, place_of_birth, date_of_birth,
	gender, phone_number, photo_url, address, npwp, grade_id, religion_id,
	echelon_id, created_at, modified_at FROM achmadnr.employees`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.Employee
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
			&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
			&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
			&employee.Address, &employee.NPWP, &employee.GradeID,
			&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
			&employee.ModifiedAt)
		if err != nil {
			return nil, err
		}
		employee.Password = "" // Clear password for security
		employees = append(employees, employee)
	}
	return employees, nil
}
func (r *EmployeeRepository) FindByID(id string) (*domain.Employee, error) {
	query := `SELECT id, role_id, nip, password, full_name, place_of_birth, date_of_birth,
	gender, phone_number, photo_url, address, npwp, grade_id, religion_id,
	echelon_id, created_at, modified_at FROM achmadnr.employees WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var employee *domain.Employee = &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
		&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
		&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
		&employee.Address, &employee.NPWP, &employee.GradeID,
		&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
		&employee.ModifiedAt)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil // Not found
		// }
		return nil, err
	}
	return employee, nil
}
func (r *EmployeeRepository) Save(employee *domain.Employee) (*domain.Employee, error) {
	query := `INSERT INTO achmadnr.employees (role_id, nip, password, full_name, place_of_birth,
	date_of_birth, gender, phone_number, photo_url, address, npwp, grade_id,
	religion_id, echelon_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
	$11, $12, $13, $14) RETURNING id, created_at, modified_at`
	err := r.db.QueryRow(query, employee.RoleID, employee.NIP, employee.Password,
		employee.FullName, employee.PlaceOfBirth, employee.DateOfBirth,
		employee.Gender, employee.PhoneNumber, employee.PhotoURL,
		employee.Address, employee.NPWP, employee.GradeID,
		employee.ReligionID, employee.EchelonID).Scan(&employee.ID, &employee.CreatedAt, &employee.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return employee, nil

}
func (r *EmployeeRepository) Update(employee *domain.Employee) (*domain.Employee, error) {
	query := `UPDATE achmadnr.employees SET role_id = $1, nip = $2, password = $3,
	full_name = $4, place_of_birth = $5, date_of_birth = $6, gender = $7,
	phone_number = $8, photo_url = $9, address = $10, npwp = $11,
	grade_id = $12, religion_id = $13, echelon_id = $14,
	modified_at = now() WHERE id = $15`
	_, err := r.db.Exec(query, employee.RoleID, employee.NIP, employee.Password,
		employee.FullName, employee.PlaceOfBirth, employee.DateOfBirth,
		employee.Gender, employee.PhoneNumber, employee.PhotoURL,
		employee.Address, employee.NPWP, employee.GradeID,
		employee.ReligionID, employee.EchelonID, employee.ID)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

// UploadProfileImage akan menyimpan url foto ke database
func (r *EmployeeRepository) UploadProfileImage(id string, fileName string) (string, error) {
	query := `UPDATE achmadnr.employees SET photo_url = $1, modified_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, fileName, id)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
func (r *EmployeeRepository) Delete(id string) error {
	query := `DELETE FROM achmadnr.employees WHERE id = $1`
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
func (r *EmployeeRepository) FindByNIP(nip string) (*domain.Employee, error) {
	query := `SELECT id, role_id, nip, password, full_name, place_of_birth, date_of_birth, gender,
	phone_number, photo_url, address, npwp, grade_id, religion_id,
	echelon_id, created_at, modified_at FROM achmadnr.employees WHERE nip = $1`
	row := r.db.QueryRow(query, nip)
	var employee *domain.Employee = &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
		&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
		&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
		&employee.Address, &employee.NPWP, &employee.GradeID,
		&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
		&employee.ModifiedAt)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil // Not found
		// }
		return nil, err
	}
	return employee, nil
}

func (r *EmployeeRepository) FindByName(name string) ([]domain.Employee, error) {
	query := `SELECT id, role_id, nip, password, full_name, place_of_birth, date_of_birth, gender,
	phone_number, photo_url, address, npwp, grade_id, religion_id,
	echelon_id, created_at, modified_at FROM achmadnr.employees WHERE full_name ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.Employee
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
			&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
			&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
			&employee.Address, &employee.NPWP, &employee.GradeID,
			&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
			&employee.ModifiedAt)
		if err != nil {
			return nil, err
		}
		employee.Password = "" // Clear password for security
		employees = append(employees, employee)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindByUnit(unitID int) ([]domain.Employee, error) {
	query := `SELECT e.id, e.role_id, e.nip, e.password, e.full_name, e.place_of_birth, e.date_of_birth, e.gender,
	e.phone_number, e.photo_url, e.address, e.npwp, e.grade_id, e.religion_id,
	e.echelon_id, e.created_at, e.modified_at
	FROM achmadnr.employees e 
	right join achmadnr.employee_assignments ea on e.id = ea.employee_id
	where ea.unit_id = $1`
	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.Employee
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
			&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
			&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
			&employee.Address, &employee.NPWP, &employee.GradeID,
			&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
			&employee.ModifiedAt)
		if err != nil {
			return nil, err
		}
		employee.Password = "" // Clear password for security
		employees = append(employees, employee)
	}
	return employees, nil
}
func (r *EmployeeRepository) Search(input string) ([]domain.Employee, error) {
	query := `SELECT id, role_id, nip, password, full_name, place_of_birth, date_of_birth, gender,
	phone_number, photo_url, address, npwp, grade_id, religion_id,
	echelon_id, created_at, modified_at FROM achmadnr.employees WHERE nip = $1 OR full_name ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(query, input)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.Employee
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(&employee.ID, &employee.RoleID, &employee.NIP, &employee.Password,
			&employee.FullName, &employee.PlaceOfBirth, &employee.DateOfBirth,
			&employee.Gender, &employee.PhoneNumber, &employee.PhotoURL,
			&employee.Address, &employee.NPWP, &employee.GradeID,
			&employee.ReligionID, &employee.EchelonID, &employee.CreatedAt,
			&employee.ModifiedAt)
		if err != nil {
			return nil, err
		}
		employee.Password = "" // Clear password for security
		employees = append(employees, employee)
	}
	return employees, nil
}
