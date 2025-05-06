package usecase_employee

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

func (eu *EmployeeUsecase) UpdateEmployee(proposerId string, nip string, employee *domain.Employee) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil || proposer == nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	// check if user is authorized to update employee or is self update are allowed
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// check wether employee already exists
	existingEmployee, err := eu.empRepo.FindByNIP(nip)
	if err != nil || existingEmployee == nil {
		return nil, &utils.InternalServerError{Message: "failed to check employee or employee not found"}
	}

	if employee.RoleID != "" {
		// unauthorized in this endpoint
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	// full name checking
	if employee.FullName != "" && len(employee.FullName) > 3 && utils.IsAlpha(employee.FullName) {
		existingEmployee.FullName = employee.FullName
	}
	// place of birth checking
	if employee.PlaceOfBirth != "" && len(employee.PlaceOfBirth) > 3 && utils.IsAlpha(employee.PlaceOfBirth) {
		existingEmployee.PlaceOfBirth = employee.PlaceOfBirth
	}
	// date of birth checking jika is Zero atau tidak ada isinya
	if !employee.DateOfBirth.IsZero() {
		existingEmployee.DateOfBirth = employee.DateOfBirth
	}
	// Gender checking
	if employee.Gender != "" && len(employee.Gender) == 1 {
		existingEmployee.Gender = employee.Gender
	}
	// phone number checking
	if employee.PhoneNumber != "" && len(employee.PhoneNumber) > 6 && utils.IsNumeric(employee.PhoneNumber) {
		existingEmployee.PhoneNumber = employee.PhoneNumber
	}
	// address checking
	if employee.Address != "" && len(employee.Address) > 6 {
		existingEmployee.Address = employee.Address
	}
	// NPWP checking
	if employee.NPWP != "" && len(employee.NPWP) == 16 {
		existingEmployee.NPWP = employee.NPWP
	}
	// grade id checking
	if employee.GradeID > 0 {
		existingEmployee.GradeID = employee.GradeID
	}
	// religion id checking
	if employee.ReligionID != "" && len(employee.ReligionID) == 3 {
		existingEmployee.ReligionID = employee.ReligionID
	}
	// echelon id checking
	if employee.EchelonID > 0 {
		existingEmployee.EchelonID = employee.EchelonID
	}

	newEmp, err := eu.empRepo.Update(existingEmployee)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to update employee"}
	}
	newEmp.Password = "" // clear password for security
	return newEmp, nil
}

func (eu *EmployeeUsecase) Promote(proposerId string, nip string, roleID string) (*domain.Employee, error) {
	// get proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	proposerRole := proposer.RoleID
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// get employee
	employee, err := eu.empRepo.FindByNIP(nip)
	employeeRole := employee.RoleID
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// Jika proposerId == employeeId maka tidak boleh!
	if proposer.ID == employee.ID {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	if employee.RoleID == roleID {
		return nil, &utils.BadRequestError{Message: "employee already has this role"}
	}
	isValid := eu.hasValidPath(proposerRole, employeeRole, roleID)
	if !isValid {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	employee.RoleID = roleID
	// save employee
	newEmployee, err := eu.empRepo.Update(employee)
	if err != nil {
		return nil, &utils.BadRequestError{Message: "role not found"}
	}
	newEmployee.Password = "" // clear password for security
	return newEmployee, nil
}
