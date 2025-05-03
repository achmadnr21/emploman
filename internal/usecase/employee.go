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
	if employee.PhotoURL == "" {
		employee.PhotoURL = "/image/profile_picture_default.png"
	}
	newEmployee, err := eu.empRepo.Save(employee)
	newEmployee.Password = "" // clear password for security
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to save employee"}
	}
	return newEmployee, nil
}

// Local Utility Functions

func validateEmployeeInput(e *domain.Employee) error {

	if e.RoleID == "" || len(e.RoleID) != 3 {
		return &utils.BadRequestError{Message: "role_id is required or invalid"}
	}
	if e.NIP == "" || len(e.NIP) != 18 || !utils.IsNumeric(e.NIP) {
		return &utils.BadRequestError{Message: "nip is required or invalid"}
	}
	if e.Password == "" || len(e.Password) < 8 {
		return &utils.BadRequestError{Message: "password is required or invalid"}
	}
	if e.FullName == "" || len(e.FullName) < 3 || !utils.IsAlpha(e.FullName) {
		return &utils.BadRequestError{Message: "full_name is required or invalid"}
	}
	if e.PlaceOfBirth == "" || len(e.PlaceOfBirth) < 3 || !utils.IsAlpha(e.PlaceOfBirth) {
		return &utils.BadRequestError{Message: "place_of_birth is required or invalid"}
	}
	if e.DateOfBirth.IsZero() {
		return &utils.BadRequestError{Message: "date_of_birth is required or invalid"}
	}
	if e.Gender == "" || len(e.Gender) != 1 {
		return &utils.BadRequestError{Message: "gender is required or invalid"}
	}
	if e.PhoneNumber == "" || len(e.PhoneNumber) < 6 || !utils.IsNumeric(e.PhoneNumber) {
		return &utils.BadRequestError{Message: "phone_number is required or invalid"}
	}
	if e.Address == "" || len(e.Address) < 6 {
		return &utils.BadRequestError{Message: "address is required or invalid"}
	}
	if e.GradeID <= 0 {
		return &utils.BadRequestError{Message: "grade_id is required or invalid"}
	}
	if e.ReligionID == "" || len(e.ReligionID) != 3 {
		return &utils.BadRequestError{Message: "religion_id is required or invalid"}
	}
	if e.EchelonID <= 0 {
		return &utils.BadRequestError{Message: "echelon_id is required or invalid"}
	}
	return nil
}
