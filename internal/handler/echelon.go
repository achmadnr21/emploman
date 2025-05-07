package handler

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type EchelonHandler struct {
	uc *usecase.EchelonUsecase
}

func NewEchelonHandler(apiV *gin.RouterGroup, uc *usecase.EchelonUsecase) {
	EchelonHandler := &EchelonHandler{
		uc: uc,
	}
	echelon := apiV.Group("/echelon")
	echelon.Use(middleware.JWTAuthMiddleware)
	{
		echelon.GET("", EchelonHandler.GetAll)
		echelon.POST("", EchelonHandler.AddEchelon)
	}
}

func (h *EchelonHandler) GetAll(c *gin.Context) {
	echelons, err := h.uc.GetAll()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError("Failed to get echelons"))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Success", echelons))
}
func (h *EchelonHandler) AddEchelon(c *gin.Context) {
	proposerId, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	var echelon domain.Echelon
	if err := c.ShouldBindJSON(&echelon); err != nil {
		c.JSON(400, utils.ResponseError("Invalid request"))
		return
	}
	if err := h.uc.AddEchelon(proposerId.(string), &echelon); err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Echelon added successfully", nil))
}
