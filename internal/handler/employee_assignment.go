package handler

import (
	"strconv"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeAssignmentHandler struct {
	uc *usecase.EmployeeAssignmentUsecase
}

func NewEmployeeAssignmentHandler(apiV *gin.RouterGroup, uc *usecase.EmployeeAssignmentUsecase) {
	empAssignHandler := &EmployeeAssignmentHandler{
		uc: uc,
	}

	empAssign := apiV.Group("/employee-assignment")
	empAssign.Use(middleware.JWTAuthMiddleware)
	{
		empAssign.GET("", empAssignHandler.GetAll)
		empAssign.POST("", empAssignHandler.AssignEmployee)
		empAssign.POST("deactivate", empAssignHandler.DeactivateEmployee)
		empAssign.GET("/:employee_id/:unit_id/:position_id", empAssignHandler.GetByAllID)
		empAssign.GET("/:employee_id", empAssignHandler.GetByEmployeeID)
	}
}
func (h *EmployeeAssignmentHandler) GetAll(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	empAssignments, err := h.uc.GetAll(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(200, empAssignments)
}

func (h *EmployeeAssignmentHandler) GetByAllID(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	employeeID := c.Param("employee_id")
	unitID := c.Param("unit_id")
	positionID := c.Param("position_id")
	unitIDInt, err := strconv.Atoi(unitID)
	if err != nil {
		c.JSON(400, utils.ResponseError("Invalid unit ID"))
		return
	}
	positionIDInt, err := strconv.Atoi(positionID)
	if err != nil {
		c.JSON(400, utils.ResponseError("Invalid position ID"))
		return
	}

	empAssign, err := h.uc.GetAssignmentByAllID(user_id.(string), employeeID, unitIDInt, positionIDInt)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(200, empAssign)
}

func (h *EmployeeAssignmentHandler) AssignEmployee(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	var empAssign *domain.EmployeeAssignment
	if err := c.ShouldBindJSON(&empAssign); err != nil {
		c.JSON(400, utils.ResponseError("Invalid request payload"))
		return
	}

	err := h.uc.AssignEmployee(user_id.(string), empAssign)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(200, utils.ResponseSuccess("Employee assigned successfully", nil))
}

func (h *EmployeeAssignmentHandler) DeactivateEmployee(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	var empAssign *domain.EmployeeAssignment
	if err := c.ShouldBindJSON(&empAssign); err != nil {
		c.JSON(400, utils.ResponseError("Invalid request payload"))
		return
	}
	err := h.uc.Deactivate(user_id.(string), empAssign.EmployeeID, empAssign.UnitID, empAssign.PositionID)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("Employee assignment deactivated successfully", nil))
}

func (h *EmployeeAssignmentHandler) GetByEmployeeID(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, utils.ResponseError("User ID not found"))
		return
	}
	employeeID := c.Param("employee_id")
	empAssign, err := h.uc.GetAssignmentByEmployeeID(user_id.(string), employeeID)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(200, utils.ResponseSuccess("", empAssign))
}
