package usecase

import (
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type EmployeeAssignmentUsecase struct {
	empAssignRepo domain.EmployeeAssignmentInterface
	empRepo       domain.EmployeeInterface
	roleRepo      domain.RoleInterface
	unitRepo      domain.UnitInterface
	positionRepo  domain.PositionInterface
}

func NewEmployeeAssignmentUsecase(empAssignRepo domain.EmployeeAssignmentInterface, empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface, positionRepo domain.PositionInterface) *EmployeeAssignmentUsecase {
	return &EmployeeAssignmentUsecase{
		empAssignRepo: empAssignRepo,
		empRepo:       empRepo,
		roleRepo:      roleRepo,
		unitRepo:      unitRepo,
		positionRepo:  positionRepo,
	}
}

func (e *EmployeeAssignmentUsecase) GetAll(proposerId string) ([]domain.EmployeeAssignmentResponse, error) {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, &utils.UnauthorizedError{Message: "Failed to get user role"}
	}
	// check proposer can assign employee global or internal
	if !proposerRole.CanAssignEmployeeGlobal && !proposerRole.CanAssignEmployeeInternal {
		return nil, &utils.UnauthorizedError{}
	}
	empAssignments, err := e.empAssignRepo.FindAll()
	if err != nil {
		fmt.Println("Error in GetAll: ", err)
		return nil, &utils.InternalServerError{Message: "Failed to get employee assignments"}
	}
	return empAssignments, nil
}

func (e *EmployeeAssignmentUsecase) GetAssignmentByAllID(proposerId string, employeeID string, unitID int, positionID int) (*domain.EmployeeAssignmentResponse, error) {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, err
	}
	// check proposer unit id via EmployeeAssignment table
	proposerAssign, err := e.empAssignRepo.FindByEmployeeID(proposerId)
	if err != nil {
		return nil, err
	}

	if !proposerRole.CanAssignEmployeeGlobal {
		if !proposerRole.CanAssignEmployeeInternal || proposerAssign.UnitID != unitID {
			return nil, &utils.UnauthorizedError{}
		}
	}

	assignmentGranted, err := e.empAssignRepo.FindByID(employeeID, unitID, positionID)
	if err != nil {
		return nil, err
	}
	return assignmentGranted, nil

}

func (e *EmployeeAssignmentUsecase) AssignEmployee(proposerId string, assignStatement *domain.EmployeeAssignment) error {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.UnauthorizedError{Message: "Failed to get user role"}
	}
	// check proposer unit id via EmployeeAssignment table
	proposerAssign, err := e.empAssignRepo.FindByEmployeeID(proposerId)
	if err != nil {
		return &utils.UnauthorizedError{Message: "Failed to get user role"}
	}

	if !proposerRole.CanAssignEmployeeGlobal {
		if !proposerRole.CanAssignEmployeeInternal || proposerAssign.UnitID != assignStatement.UnitID {
			return &utils.UnauthorizedError{}
		}
	}
	// check existance of employee
	emp, _ := e.empRepo.FindByID(assignStatement.EmployeeID)
	if emp == nil {
		return &utils.NotFoundError{Message: "Employee not found"}
	}
	// check existance of unit
	unit, _ := e.unitRepo.FindByID(assignStatement.UnitID)
	if unit == nil {
		return &utils.NotFoundError{Message: "Unit not found"}
	}
	// check existance of position
	position, _ := e.positionRepo.FindByID(assignStatement.PositionID)
	if position == nil {
		return &utils.NotFoundError{Message: "Position not found"}
	}
	assignStatement.IsActive = true
	// perform transactional assignment
	err = e.empAssignRepo.TransactionalAssignment(assignStatement)
	if err != nil {
		fmt.Println("Error in AssignEmployee: ", err)
		return &utils.InternalServerError{Message: "Failed to assign employee"}
	}

	return nil
}

func (e *EmployeeAssignmentUsecase) Deactivate(proposerId string, employeeID string, unitID int, positionID int) error {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.UnauthorizedError{Message: "Failed to get user role"}
	}
	// check proposer unit id via EmployeeAssignment table
	proposerAssign, err := e.empAssignRepo.FindByEmployeeID(proposerId)
	if err != nil {
		return &utils.UnauthorizedError{Message: "Failed to get user role"}
	}

	if !proposerRole.CanAssignEmployeeGlobal {
		if !proposerRole.CanAssignEmployeeInternal || proposerAssign.UnitID != unitID {
			return &utils.UnauthorizedError{Message: "Unauthorized to deactivate employee assignment"}
		}
	}

	err = e.empAssignRepo.Deactivate(employeeID, unitID, positionID)
	if err != nil {
		fmt.Println("Error in Deactivate: ", err)
		return &utils.InternalServerError{Message: "Failed to deactivate employee assignment"}
	}
	return nil
}

func (e *EmployeeAssignmentUsecase) GetAssignmentByEmployeeID(proposerId string, employeeID string) (*domain.EmployeeAssignmentResponse, error) {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, &utils.UnauthorizedError{Message: "Failed to get user role"}
	}
	// check proposer unit id via EmployeeAssignment table
	proposerAssign, err := e.empAssignRepo.FindByEmployeeID(proposerId)
	if err != nil {
		return nil, &utils.UnauthorizedError{Message: "Failed to get user role"}
	}

	if !proposerRole.CanAssignEmployeeGlobal {
		if !proposerRole.CanAssignEmployeeInternal || proposerAssign.EmployeeID != employeeID {
			return nil, &utils.UnauthorizedError{}
		}
	}

	assignments, err := e.empAssignRepo.FindByEmployeeID(employeeID)
	if err != nil {
		fmt.Println("Error in GetAssignmentByEmployeeID: ", err)
		return nil, &utils.InternalServerError{Message: "Failed to get employee assignment"}
	}
	return assignments, nil
}
