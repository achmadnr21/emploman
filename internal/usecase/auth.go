package usecase

import (
	"time"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
)

type AuthUsecase struct {
	EmpRepo  domain.EmployeeInterface
	RoleRepo domain.RoleInterface
}

func NewAuthUsecase(employeeRepo domain.EmployeeInterface, roleRepo domain.RoleInterface) *AuthUsecase {
	return &AuthUsecase{
		EmpRepo:  employeeRepo,
		RoleRepo: roleRepo,
	}
}

func (au *AuthUsecase) Login(nip, password string) (string, string, error) {
	// Find employee by NIK
	employee, err := au.EmpRepo.FindByNIP(nip)
	if err != nil {
		return "", "", &utils.NotFoundError{Message: "user not found"}
	}
	// Check password
	if !utils.CheckPasswordHash(password, employee.Password) {
		return "", "", &utils.UnauthorizedError{Message: "invalid user or password"}
	}
	// generate token
	token, err := utils.GenerateAccessToken(employee.ID)
	if err != nil {
		return "", "", &utils.InternalServerError{Message: "failed to generate token"}
	}
	refreshToken, err := utils.GenerateRefreshToken(employee.ID)
	if err != nil {
		return "", "", &utils.InternalServerError{Message: "failed to generate refresh token"}
	}
	// Return token and role
	return token, refreshToken, nil

}

func (au *AuthUsecase) RefreshToken(refreshToken string) (string, string, error) {
	// Parse refresh token
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", &utils.UnauthorizedError{Message: "invalid refresh token"}
	}
	// expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", "", &utils.UnauthorizedError{Message: "refresh token expired"}
	}
	// generate token ketika valid dan tidak expired
	token, err := utils.GenerateAccessToken(claims.UserId)
	if err != nil {
		return "", "", &utils.InternalServerError{Message: "failed to generate token"}
	}
	refreshToken, err = utils.GenerateRefreshToken(claims.UserId)
	if err != nil {
		return "", "", &utils.InternalServerError{Message: "failed to generate refresh token"}
	}
	// Return token and role
	return token, refreshToken, nil
}
