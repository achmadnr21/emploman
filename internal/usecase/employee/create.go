package usecase_employee

import (
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

func (eu *EmployeeUsecase) Add(proposerId string, employee *domain.Employee) (*domain.Employee, error) {
	// check employee.RoleID should be empty
	if employee.RoleID != "" {
		return nil, &utils.BadRequestError{Message: "Invalid Payload"}
	}
	employee.RoleID = "USR"
	if _, _, err := eu.authorize(proposerId, true); err != nil {
		return nil, err
	}
	// validate employee input
	if err := validateEmployeeInput(employee); err != nil {
		return nil, err
	}
	// check wether employee already exists
	existingEmployee, err := eu.empRepo.FindByNIP(employee.NIP)
	if err == nil && existingEmployee == nil {
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
	employee.PhotoURL = "https://s3.nevaobjects.id/emploman/pictureprofile/defaultprofile.jpg"
	// save employee
	newEmployee, err := eu.empRepo.Save(employee)
	// newEmployee.Password = "" // clear password for security
	if err != nil {
		fmt.Println("Error saving employee:", err)
		return nil, &utils.InternalServerError{Message: "failed to save employee"}
	}
	newEmployee.Password = ""
	return newEmployee, nil
}
