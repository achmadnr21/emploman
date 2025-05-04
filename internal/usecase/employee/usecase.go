package usecase_employee

import (
	"github.com/achmadnr21/emploman/internal/domain"
)

type EmployeeUsecase struct {
	empRepo  domain.EmployeeInterface
	roleRepo domain.RoleInterface
	unitRepo domain.UnitInterface
	s3Repo   domain.S3Interface
}

func NewEmployeeUsecase(empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface, s3Repo domain.S3Interface) *EmployeeUsecase {
	return &EmployeeUsecase{
		empRepo:  empRepo,
		roleRepo: roleRepo,
		unitRepo: unitRepo,
		s3Repo:   s3Repo,
	}
}
