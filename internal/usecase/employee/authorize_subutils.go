package usecase_employee

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

func (eu *EmployeeUsecase) authorize(proposerId string, requireAdd bool) (*domain.Employee, *domain.Role, error) {
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

func (eu *EmployeeUsecase) hasValidPath(proposerRole string, employeeRole string, roleID string) bool {
	// check using uc
	promoteList, err := eu.roleRepo.FindPromoteRole(proposerRole)
	if err != nil {
		return false
	}

	return validateRolePromotion(employeeRole, roleID, promoteList)
}

// ==================================================================== UTILITIES ====================================================================

func validateRolePromotion(currRole string, targetRole string, roleList []domain.RolePromotion) bool {
	var currExists bool = false
	var targetExists bool = false
	for _, role := range roleList {
		// curr exists on the list
		if currRole == role.FromRoleID || currRole == role.ToRoleID {
			currExists = true
		}
		if targetRole == role.FromRoleID || targetRole == role.ToRoleID {
			targetExists = true
		}

	}
	return currExists && targetExists
}
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
