package usecase

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type UnitUsecase struct {
	unitRepo domain.UnitInterface
	roleRepo domain.RoleInterface
}

func NewUnitUsecase(unitRepo domain.UnitInterface, roleRepo domain.RoleInterface) *UnitUsecase {
	return &UnitUsecase{
		unitRepo: unitRepo,
		roleRepo: roleRepo,
	}
}

func (uc *UnitUsecase) GetAllUnit() ([]domain.Unit, error) {
	units, err := uc.unitRepo.FindAll()
	if err != nil {
		return nil, &utils.NotFoundError{Message: "unit not found"}
	}
	return units, nil
}
func (uc *UnitUsecase) AddUnit(proposerId string, unit *domain.Unit) (*domain.Unit, error) {
	// cek proposer role.
	proposer, err := uc.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// check if proposer can add unit
	if !proposer.CanAddUnit {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to add unit"}
	}
	// checking all fillable
	if unit.Name == "" || len(unit.Name) < 5 {
		return nil, &utils.BadRequestError{Message: "unit name is required and must be at least 5 characters"}
	}
	if unit.Address == "" || len(unit.Address) < 10 {
		return nil, &utils.BadRequestError{Message: "unit address is required and must be at least 10 characters"}
	}
	if unit.Description == "" || len(unit.Description) < 5 {
		return nil, &utils.BadRequestError{Message: "unit description is required and must be at least 5 characters"}
	}
	newunit, err := uc.unitRepo.Save(unit)
	if err != nil {
		return nil, err
	}
	return newunit, nil
}

func (uc *UnitUsecase) UpdateUnit(proposerId string, unit *domain.Unit) (*domain.Unit, error) {
	// cek proposer role.

	proposer, err := uc.roleRepo.FindByUserID(proposerId)

	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// check if proposer can add unit
	if !proposer.CanAddUnit {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to add unit"}
	}

	// get unit by id
	oldunit, err := uc.unitRepo.FindByID(unit.ID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "unit not found"}
	}
	// checking all fillable
	if unit.Name != "" && len(unit.Name) > 5 {
		oldunit.Name = unit.Name
	}
	if unit.Description != "" && len(unit.Description) > 5 {
		oldunit.Description = unit.Description
	}
	if unit.Address != "" && len(unit.Address) > 10 {
		oldunit.Address = unit.Address
	}

	newunit, err := uc.unitRepo.Update(oldunit)
	if err != nil {
		return nil, err
	}
	return newunit, nil
}
func (uc *UnitUsecase) DeleteUnit(proposerId string, id int) error {
	// cek proposer role.
	proposer, err := uc.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.NotFoundError{Message: "user not found"}
	}
	// check if proposer can add unit
	if !proposer.CanAddUnit {
		return &utils.UnauthorizedError{Message: "user not authorized to add unit"}
	}
	err = uc.unitRepo.Delete(id)
	if err != nil {
		return &utils.NotFoundError{Message: "unit not found"}
	}
	return nil
}

func (uc *UnitUsecase) GetUnitByID(id int) (*domain.Unit, error) {
	unit, err := uc.unitRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (uc *UnitUsecase) SearchUnit(input string) ([]domain.Unit, error) {
	units, err := uc.unitRepo.Search(input)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "unit not found"}
	}
	return units, nil
}
