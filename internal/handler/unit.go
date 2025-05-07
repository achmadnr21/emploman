package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	uc *usecase.UnitUsecase
}

func NewUnitHandler(apiV *gin.RouterGroup, uc *usecase.UnitUsecase) {
	unitHandler := &UnitHandler{
		uc: uc,
	}

	unit := apiV.Group("/unit")
	unit.Use(middleware.JWTAuthMiddleware)
	{
		unit.GET("", unitHandler.GetAllUnit) // GET /units
		unit.POST("", unitHandler.AddUnit)   // POST /units
		unit.GET("/:id", unitHandler.GetUnitByID)
		unit.PUT("/:id", unitHandler.UpdateUnit)
		unit.DELETE("/:id", unitHandler.DeleteUnit)
		unit.GET("/search", unitHandler.SearchUnit) // GET /units/search
	}
}
func (h *UnitHandler) GetAllUnit(c *gin.Context) {
	units, err := h.uc.GetAllUnit()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get all unit", units))
}
func (h *UnitHandler) AddUnit(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var payload domain.Unit
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("error bind json", err)
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	unit, err := h.uc.AddUnit(userId.(string), &payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Add unit", unit))
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	userId, _ := c.Get("user_id")
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}
	var payload domain.Unit
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	payload.ID = idInt
	unit, err := h.uc.UpdateUnit(userId.(string), &payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Update unit", unit))
}
func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	userId, _ := c.Get("user_id")
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}
	err = h.uc.DeleteUnit(userId.(string), idInt)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Delete unit", nil))
}
func (h *UnitHandler) SearchUnit(c *gin.Context) {
	var payload struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	units, err := h.uc.SearchUnit(payload.Query)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Search unit", units))
}

func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id := c.Param("id")
	// Convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}
	unit, err := h.uc.GetUnitByID(idInt)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get unit by ID", unit))
}
