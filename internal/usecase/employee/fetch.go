package usecase_employee

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

func (eu *EmployeeUsecase) GetAll(proposerId string) ([]domain.Employee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
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
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
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
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
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
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get employee by input
	employees, err := eu.empRepo.Search(input)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}
