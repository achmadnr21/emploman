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

/*
--
--employees ddl

create table achmadnr.employees(
	id uuid unique default gen_random_uuid() primary key,
	role_id char(3) not null,
	nip varchar(20) unique not null,
	password varchar(255) not null,
	full_name varchar(255) not null,
	place_of_birth varchar(100) not null,
	date_of_birth DATE not null,
	gender char(1) check (gender in('L','P')),
	phone_number varchar(20) not null,
	photo_url text default '/image/profile_picture_default.png',
	address text not null,
	npwp varchar(25) null,
	grade_id int not null,
	religion_id char(3) not null,
	echelon_id int not null,
	created_at timestamp default now(),
	modified_at timestamp default now(),
	foreign key (grade_id) references achmadnr.grades(id) on delete set null,
	foreign key (religion_id) references achmadnr.religions(id) on delete set null,
	foreign key (echelon_id) references achmadnr.echelons(id) on delete set null,
	foreign key (role_id) references achmadnr.roles(id) on delete set null
);


type Employee struct {
	ID           string    `json:"id" db:"id"`
	RoleID       string    `json:"role_id" db:"role_id"`
	NIP          string    `json:"nip" db:"nip"`
	Password     string    `json:"password" db:"password"`
	FullName     string    `json:"full_name" db:"full_name"`
	PlaceOfBirth string    `json:"place_of_birth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender       string    `json:"gender" db:"gender"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	PhotoURL     string    `json:"photo_url" db:"photo_url"`
	Address      string    `json:"address" db:"address"`
	NPWP         *string   `json:"npwp" db:"npwp"`
	GradeID      int       `json:"grade_id" db:"grade_id"`
	ReligionID   string    `json:"religion_id" db:"religion_id"`
	EchelonID    int       `json:"echelon_id" db:"echelon_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ModifiedAt   time.Time `json:"modified_at" db:"modified_at"`
}

type EmployeeInterface interface {
	FindAll() ([]Employee, error)
	FindByID(id string) (*Employee, error)
	Save(employee *Employee) (*Employee, error)
	Update(employee *Employee) (*Employee, error)
	UploadProfileImage(id string, fileName string) (string, error)
	Delete(id string) error
	FindByNIP(nip string) (*Employee, error)
	FindByName(name string) ([]Employee, error)
}

*/

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
