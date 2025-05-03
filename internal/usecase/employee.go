package usecase

import (
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type EmployeeUsecase struct {
	empRepo  domain.EmployeeInterface
	roleRepo domain.RoleInterface
	unitRepo domain.UnitInterface
}

func NewEmployeeUsecase(empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface) *EmployeeUsecase {
	return &EmployeeUsecase{
		empRepo:  empRepo,
		roleRepo: roleRepo,
		unitRepo: unitRepo,
	}
}

func (eu *EmployeeUsecase) GetAll(proposerId string) ([]domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// get all employee
	employees, err := eu.empRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get employees"}
	}
	// return employees
	return employees, nil
}

func (eu *EmployeeUsecase) GetByNIP(proposerId string, nip string) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// get employee by nip
	employee, err := eu.empRepo.FindByNIP(nip)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employee, nil
}
func (eu *EmployeeUsecase) GetByUnit(proposerId string, unitId int) ([]domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// check wether unit exists
	unit, err := eu.unitRepo.FindByID(unitId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "unit not found"}
	}

	// get employee by unit
	employees, err := eu.empRepo.FindByUnit(unit.ID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}

func (eu *EmployeeUsecase) Search(proposerId string, input string) ([]domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// get employee by input
	employees, err := eu.empRepo.Search(input)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}

/*
package domain
import (
	"time"
)
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
*/

/*
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
*/

func validateEmployeeInput(e *domain.Employee) error {

	if e.RoleID == "" {
		return &utils.BadRequestError{Message: "role_id is required"}
	}
	if e.NIP == "" {
		return &utils.BadRequestError{Message: "nip is required"}
	}
	if e.Password == "" {
		return &utils.BadRequestError{Message: "password is required"}
	}
	if e.FullName == "" {
		return &utils.BadRequestError{Message: "full_name is required"}
	}
	if e.PlaceOfBirth == "" {
		return &utils.BadRequestError{Message: "place_of_birth is required"}
	}
	if e.DateOfBirth.IsZero() {
		return &utils.BadRequestError{Message: "date_of_birth is required"}
	}
	if e.Gender == "" {
		return &utils.BadRequestError{Message: "gender is required"}
	}
	if e.PhoneNumber == "" {
		return &utils.BadRequestError{Message: "phone_number is required"}
	}
	if e.Address == "" {
		return &utils.BadRequestError{Message: "address is required"}
	}
	if e.GradeID <= 0 {
		return &utils.BadRequestError{Message: "grade_id is required"}
	}
	if e.ReligionID == "" {
		return &utils.BadRequestError{Message: "religion_id is required"}
	}
	if e.EchelonID <= 0 {
		return &utils.BadRequestError{Message: "echelon_id is required"}
	}
	return nil
}

func (eu *EmployeeUsecase) Add(proposerId string, employee *domain.Employee) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if !pr.CanAddEmployee {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// validate employee input
	if err := validateEmployeeInput(employee); err != nil {
		return nil, err
	}
	// check wether employee already exists
	existingEmployee, err := eu.empRepo.FindByNIP(employee.NIP)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to check employee"}
	}
	if existingEmployee != nil {
		return nil, &utils.ConflictError{Message: "employee already exists"}
	}

	// now, hash the password
	hashedPassword, err := utils.HashPassword(employee.Password)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to hash password"}
	}
	employee.Password = hashedPassword
	// save employee
	newEmployee, err := eu.empRepo.Save(employee)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to save employee"}
	}
	return newEmployee, nil
}
