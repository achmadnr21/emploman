package usecase

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type GradeUsecase struct {
	GradeRepo domain.GradeInterface
	roleRepo  domain.RoleInterface
}

func NewGradeUsecase(gradeRepo domain.GradeInterface, roleRepo domain.RoleInterface) *GradeUsecase {
	return &GradeUsecase{
		GradeRepo: gradeRepo,
		roleRepo:  roleRepo,
	}
}

func (g *GradeUsecase) GetAll() ([]domain.Grade, error) {
	grades, err := g.GradeRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "Failed to get grades"}
	}
	return grades, nil
}
func (g *GradeUsecase) AddGrade(proposerId string, grade *domain.Grade) error {
	// check proposer role
	proposerRole, err := g.roleRepo.FindByUserID(proposerId)
	if err != nil {
		return &utils.InternalServerError{Message: "Failed to get proposer role"}
	}
	// check proposer can add grade
	if !proposerRole.CanAddGrade {
		return &utils.UnauthorizedError{Message: "You cannot add grade"}
	}

	// make sure grade code is not empty
	if grade.Code == "" {
		return &utils.BadRequestError{Message: "Grade code cannot be empty"}
	}
	if _, err := g.GradeRepo.Save(grade); err != nil {
		return &utils.InternalServerError{Message: "Failed to add grade possibly duplicate ID"}
	}
	return nil
}
