package usecase

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type EchelonUsecase struct {
	EchelonRepo domain.EchelonInterface
	roleRepo    domain.RoleInterface
}

func NewEchelonUsecase(echelonRepo domain.EchelonInterface, roleRepo domain.RoleInterface) *EchelonUsecase {
	return &EchelonUsecase{
		EchelonRepo: echelonRepo,
		roleRepo:    roleRepo,
	}
}
func (e *EchelonUsecase) GetAll() ([]domain.Echelon, error) {
	echelons, err := e.EchelonRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "Failed to get echelons"}
	}
	return echelons, nil
}
func (e *EchelonUsecase) AddEchelon(proposerId string, echelon *domain.Echelon) error {
	// check proposer role
	proposerRole, err := e.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.InternalServerError{Message: "Failed to get proposer role"}
	}
	// check proposer can add echelon
	if !proposerRole.CanAddEchelon {
		return &utils.UnauthorizedError{Message: "You cannot add echelon"}
	}

	// make sure echelon code is not empty
	if echelon.Code == "" {
		return &utils.BadRequestError{Message: "Echelon code cannot be empty"}
	}
	if _, err := e.EchelonRepo.Save(echelon); err != nil {
		return &utils.InternalServerError{Message: "Failed to add echelon possibly duplicate ID"}
	}
	return nil
}
