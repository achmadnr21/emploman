package usecase

import (
	"fmt"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type ReligionUsecase struct {
	ReligionRepo domain.ReligionInterface
	roleRepo     domain.RoleInterface
}

func NewReligionUsecase(religionRepo domain.ReligionInterface, roleRepo domain.RoleInterface) *ReligionUsecase {
	return &ReligionUsecase{
		ReligionRepo: religionRepo,
		roleRepo:     roleRepo,
	}
}
func (r *ReligionUsecase) GetAll() ([]domain.Religion, error) {
	religions, err := r.ReligionRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "Failed to get religions"}
	}
	return religions, nil
}
func (r *ReligionUsecase) AddReligion(proposerId string, religion *domain.Religion) error {
	// check proposer role
	proposerRole, err := r.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.UnauthorizedError{Message: "Failed to get user role"}
	}
	// check proposer can add religion
	if !proposerRole.CanAddReligion {
		return &utils.UnauthorizedError{Message: "You cannot add religion"}
	}
	// make sure religio.ID is char(3)
	if len(religion.ID) != 3 {
		return &utils.BadRequestError{Message: "Religion ID must be 3 characters"}
	}
	// make sure name is not empty
	if religion.Name == "" || len(religion.Name) < 3 {
		return &utils.BadRequestError{Message: "Religion name cannot be empty"}
	}

	if _, err := r.ReligionRepo.Save(religion); err != nil {
		return &utils.InternalServerError{Message: fmt.Sprintf("Failed to add religion possibly duplicate ID")}
	}
	return nil
}
