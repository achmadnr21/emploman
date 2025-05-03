package handler

import (
	"net/http"
	"strings"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}
func (h *AuthHandler) Login(c *gin.Context) {
	var employee domain.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}

	token, refreshToken, err := h.uc.Login(employee.NIP, employee.Password)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	var loginResponse = struct {
		Token   string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}{
		Token:   token,
		Refresh: refreshToken,
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Login successful", loginResponse))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		c.Abort()
		return
	}
	token, refreshToken, err := h.uc.RefreshToken(tokenString)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	var refreshResponse = struct {
		Token   string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}{
		Token:   token,
		Refresh: refreshToken,
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Token refreshed", refreshResponse))
}
