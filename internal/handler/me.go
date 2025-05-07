package handler

import (
	"net/http"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	uc *usecase.MeUsecase
}

func NewMeHandler(apiV *gin.RouterGroup, uc *usecase.MeUsecase) {
	meHandler := &MeHandler{
		uc: uc,
	}

	me := apiV.Group("/me")
	me.Use(middleware.JWTAuthMiddleware)
	{
		me.GET("", meHandler.GetMe)
		me.PUT("", meHandler.UpdateMe)
		me.POST("/profile-picture", meHandler.UploadPPMe)
	}
}

func (h *MeHandler) GetMe(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	employee, err := h.uc.GetMe(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get employee by NIP", employee))
}
func (h *MeHandler) UpdateMe(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	var payload domain.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employee, err := h.uc.UpdateMe(user_id.(string), &payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Update employee", employee))
}

func (h *MeHandler) UploadPPMe(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid file"))
		return
	}
	url, err := h.uc.UploadPPMe(user_id.(string), file)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	var payload struct {
		URL string `json:"url"`
	}
	payload.URL = url
	c.JSON(http.StatusOK, utils.ResponseSuccess("Upload profile picture", payload))
}
