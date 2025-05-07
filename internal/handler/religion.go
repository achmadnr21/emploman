package handler

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type ReligionHandler struct {
	uc *usecase.ReligionUsecase
}

func NewReligionHandler(apiV *gin.RouterGroup, uc *usecase.ReligionUsecase) {
	religionHandler := &ReligionHandler{
		uc: uc,
	}

	religion := apiV.Group("/religion")
	religion.Use(middleware.JWTAuthMiddleware)
	{
		religion.GET("", religionHandler.GetAll)
		religion.POST("", religionHandler.AddReligion)
	}
}
func (h *ReligionHandler) GetAll(c *gin.Context) {
	religions, err := h.uc.GetAll()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError("Failed to get religions"))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Success", religions))
}
func (h *ReligionHandler) AddReligion(c *gin.Context) {
	proposerId, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	var religion domain.Religion
	if err := c.ShouldBindJSON(&religion); err != nil {
		c.JSON(400, utils.ResponseError("Invalid request"))
		return
	}
	if err := h.uc.AddReligion(proposerId.(string), &religion); err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Religion added successfully", nil))
}
