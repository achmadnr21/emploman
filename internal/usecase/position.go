package usecase

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type PositionUsecase struct {
	positionRepo domain.PositionInterface
	roleRepo     domain.RoleInterface
}

func NewPositionUsecase(positionRepo domain.PositionInterface, roleRepo domain.RoleInterface) *PositionUsecase {
	return &PositionUsecase{
		positionRepo: positionRepo,
		roleRepo:     roleRepo,
	}
}

func (uc *PositionUsecase) AddPosition(proposerId string, position *domain.Position) (*domain.Position, error) {
	// get proposer role
	proposerRole, err := uc.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, err
	}
	// jika tidak bisa add position, maka return error
	if !proposerRole.CanAddPosition {
		return nil, &utils.UnauthorizedError{Message: "You are not authorized to add position"}
	}
	// proses position
	position, err = uc.positionRepo.Save(position)
	if err != nil {
		return nil, &utils.InternalServerError{Message: err.Error()}
	}
	return position, nil
}
func (uc *PositionUsecase) UpdatePosition(proposerId string, position *domain.Position) (*domain.Position, error) {
	// get proposer role
	proposerRole, err := uc.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, err
	}
	// jika tidak bisa add position, maka return error
	if !proposerRole.CanAddPosition {
		return nil, &utils.UnauthorizedError{Message: "You are not authorized to update position"}
	}
	// check if position exists
	oldPosition, err := uc.positionRepo.FindByID(position.ID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "Position not found"}
	}

	// perform all checking
	// check if position name already exists
	positions, err := uc.positionRepo.FindByName(position.Name)
	if err != nil {
		return nil, err
	}
	if len(positions) > 0 {
		for _, p := range positions {
			if p.ID != oldPosition.ID {
				return nil, &utils.ConflictError{Message: "Position name already exists"}
			}
		}
	}
	// check if position id is valid
	if position.ID <= 0 {
		return nil, &utils.BadRequestError{Message: "Invalid position id"}
	}
	// check if position name is not empty and has more than 5 characters
	if position.Name != "" && len(position.Name) > 5 {
		oldPosition.Name = position.Name
	}

	// proses position
	position, err = uc.positionRepo.Update(oldPosition)
	if err != nil {
		return nil, err
	}
	return position, nil
}
func (uc *PositionUsecase) DeletePosition(proposerId string, id int) error {
	// get proposer role
	proposerRole, err := uc.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return err
	}
	// jika tidak bisa delete position, maka return error
	if !proposerRole.CanAddPosition {
		return &utils.UnauthorizedError{Message: "You are not authorized to delete position"}
	}
	// proses position
	err = uc.positionRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (uc *PositionUsecase) GetAllPosition() ([]domain.Position, error) {
	positions, err := uc.positionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return positions, nil
}
func (uc *PositionUsecase) GetPositionByID(id int) (*domain.Position, error) {
	position, err := uc.positionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return position, nil
}
func (uc *PositionUsecase) SearchPosition(query string) ([]domain.Position, error) {
	positions, err := uc.positionRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return positions, nil
}
