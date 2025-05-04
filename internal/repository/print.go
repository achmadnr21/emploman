package repository

import (
	"database/sql"

	"github.com/achmadnr21/emploman/internal/domain"
)

type PrintRepository struct {
	db *sql.DB
}

func NewPrintRepository(db *sql.DB) *PrintRepository {
	return &PrintRepository{
		db: db,
	}
}

func (r *PrintRepository) PrintAll() ([]domain.PrintEmployee, error) {
	query := `select 
		ae.nip, ae.full_name, ae.place_of_birth, ae.address, ae.date_of_birth, ae.gender, ag.code, aec.code,
		COALESCE(ap.name, '-') as position_name,
		coalesce(au.address, '-') as tempat_kerja,
		ar.name,
		COALESCE(au.name, '-') AS unit_name,
		ae.phone_number,  ae.photo_url, COALESCE(ae.npwp, '-') as npwp
		from achmadnr.employees ae
		left join achmadnr.employee_assignments aea on ae.id = aea.employee_id
		left join achmadnr.units au on aea.unit_id = au.id
		left join achmadnr.positions ap on aea.position_id = ap.id
		left join achmadnr.grades ag on ae.grade_id = ag.id
		left join achmadnr.echelons aec on ae.echelon_id = aec.id
		left join achmadnr.religions ar on ae.religion_id = ar.id;
		`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.PrintEmployee
	for rows.Next() {
		var employee domain.PrintEmployee
		err := rows.Scan(&employee.NIP, &employee.FullName, &employee.PlaceOfBirth, &employee.Address,
			&employee.DateOfBirth, &employee.Gender, &employee.Grade, &employee.Echelon,
			&employee.Jabatan, &employee.TempatTugas, &employee.Religion, &employee.Unit,
			&employee.PhoneNumber, &employee.PhotoURL, &employee.NPWP)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}
func (r *PrintRepository) PrintByNIP(id string) (*domain.PrintEmployee, error) {
	query := `select
		ae.nip, ae.full_name, ae.place_of_birth, ae.address, ae.date_of_birth,
		gender, ag.code, aec.code,
		COALESCE(ap.name, '-') as position_name,
		coalesce(au.address, '-') as tempat_kerja,
		ar.name,
		COALESCE(au.name, '-') AS unit_name,
		ae.phone_number,  ae.photo_url, COALESCE(ae.npwp, '-') as npwp
		from achmadnr.employees ae
		left join achmadnr.employee_assignments aea on ae.id = aea.employee_id
		left join achmadnr.units au on aea.unit_id = au.id
		left join achmadnr.positions ap on aea.position_id = ap.id
		left join achmadnr.grades ag on ae.grade_id = ag.id
		left join achmadnr.echelons aec on ae.echelon_id = aec.id
		left join achmadnr.religions ar on ae.religion_id = ar.id
		where ae.nip = $1`
	row := r.db.QueryRow(query, id)
	var employee *domain.PrintEmployee = &domain.PrintEmployee{}
	err := row.Scan(&employee.NIP, &employee.FullName, &employee.PlaceOfBirth, &employee.Address,
		&employee.DateOfBirth, &employee.Gender, &employee.Grade, &employee.Echelon,
		&employee.Jabatan, &employee.TempatTugas, &employee.Religion, &employee.Unit,
		&employee.PhoneNumber, &employee.PhotoURL, &employee.NPWP)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil // Not found
		// }
		return nil, err
	}
	return employee, nil
}
func (r *PrintRepository) PrintByUnit(unitID int) ([]domain.PrintEmployee, error) {
	query := `select
		ae.nip, ae.full_name, ae.place_of_birth, ae.address, ae.date_of_birth, ae.gender, ag.code, aec.code,
		COALESCE(ap.name, '-') as position_name,
		coalesce(au.address, '-') as tempat_kerja,
		ar.name,
		COALESCE(au.name, '-') AS unit_name,
		ae.phone_number,  ae.photo_url, COALESCE(ae.npwp, '-') as npwp
		from achmadnr.employees ae
		left join achmadnr.employee_assignments aea on ae.id = aea.employee_id
		left join achmadnr.units au on aea.unit_id = au.id
		left join achmadnr.positions ap on aea.position_id = ap.id
		left join achmadnr.grades ag on ae.grade_id = ag.id
		left join achmadnr.echelons aec on ae.echelon_id = aec.id
		left join achmadnr.religions ar on ae.religion_id = ar.id
		where aea.unit_id = $1`
	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []domain.PrintEmployee
	for rows.Next() {
		var employee domain.PrintEmployee
		err := rows.Scan(&employee.NIP, &employee.FullName, &employee.PlaceOfBirth, &employee.Address,
			&employee.DateOfBirth, &employee.Gender, &employee.Grade, &employee.Echelon,
			&employee.Jabatan, &employee.TempatTugas, &employee.Religion, &employee.Unit,
			&employee.PhoneNumber, &employee.PhotoURL, &employee.NPWP)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}
