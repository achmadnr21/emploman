package usecase

import (
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type PrintUsecase struct {
	printRepo domain.PrintInterface
	empRepo   domain.EmployeeInterface
	roleRepo  domain.RoleInterface
	unitRepo  domain.UnitInterface
}

func NewPrintUsecase(printRepo domain.PrintInterface, empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface) *PrintUsecase {
	return &PrintUsecase{
		printRepo: printRepo,
		empRepo:   empRepo,
		roleRepo:  roleRepo,
		unitRepo:  unitRepo,
	}
}

// private:
func (eu *PrintUsecase) authorize(proposerId string, requireAdd bool) (*domain.Employee, *domain.Role, error) {
	user, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, nil, &utils.NotFoundError{Message: "user not found"}
	}
	role, err := eu.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if requireAdd && !role.CanAddEmployee {
		return nil, nil, &utils.UnauthorizedError{Message: "not authorized"}
	}
	if !requireAdd && !(role.CanAddEmployee || role.CanAssignEmployeeInternal || role.CanAssignEmployeeGlobal) {
		return nil, nil, &utils.UnauthorizedError{Message: "not authorized"}
	}
	return user, role, nil
}

func (eu *PrintUsecase) PrintAll(proposerId string) ([]domain.PrintEmployee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get all employee
	employees, err := eu.printRepo.PrintAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get employees"}
	}
	// return employees
	return employees, nil
}

func (eu *PrintUsecase) PrintByUnitID(proposerId string, unitId int) ([]domain.PrintEmployee, error) {
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
	employees, err := eu.printRepo.PrintByUnit(unit.ID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}

func (eu *PrintUsecase) PrintByNIP(proposerId string, nip string) (*domain.PrintEmployee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get employee by nip
	employee, err := eu.printRepo.PrintByNIP(nip)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employee, nil
}
