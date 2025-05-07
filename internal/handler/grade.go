package handler

import (
	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	uc *usecase.GradeUsecase
}

func NewGradeHandler(apiV *gin.RouterGroup, uc *usecase.GradeUsecase) {
	gradeHandler := &GradeHandler{
		uc: uc,
	}

	grade := apiV.Group("/grade")
	grade.Use(middleware.JWTAuthMiddleware)
	{
		grade.GET("", gradeHandler.GetAll)
		grade.POST("", gradeHandler.AddGrade)
	}
}
func (h *GradeHandler) GetAll(c *gin.Context) {
	grades, err := h.uc.GetAll()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError("Failed to get grades"))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Success", grades))
}
func (h *GradeHandler) AddGrade(c *gin.Context) {
	proposerId, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	var grade domain.Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(400, utils.ResponseError("Invalid request"))
		return
	}
	if err := h.uc.AddGrade(proposerId.(string), &grade); err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Grade added successfully", nil))
}
