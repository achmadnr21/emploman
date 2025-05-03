/*
// 2. Employee
	employee := apiV.Group("/employee")
	employee.Use(middleware.JWTAuthMiddleware)
	{
		employee.GET("", empHandler.GetAll)
		employee.GET("/:nip", empHandler.GetByNIP)
		employee.GET("/:unit_id", empHandler.GetByUnit)
		employee.GET("/search", empHandler.Search)
		employee.POST("/add", empHandler.Add)
		employee.PUT("/update", empHandler.UpdateEmployee)
	}
*/

package handler

import (
	"net/http"
	"strconv"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	uc *usecase.EmployeeUsecase
}

func NewEmployeeHandler(uc *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		uc: uc,
	}
}

/*
employee.Use(middleware.JWTAuthMiddleware)

	{
		employee.GET("", empHandler.GetAll)
		employee.GET("/:nip", empHandler.GetByNIP)
		employee.GET("/:unit_id", empHandler.GetByUnit)
		employee.GET("/search/:name", empHandler.SearchByName)
		employee.POST("/add", empHandler.Add)
		employee.PUT("/update", empHandler.UpdateEmployee)
	}
*/
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	employees, err := h.uc.GetAll(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get all employees", employees))
}

func (h *EmployeeHandler) GetByNIP(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	employee, err := h.uc.GetByNIP(user_id.(string), nip)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get employee by NIP", employee))
}
func (h *EmployeeHandler) GetByUnit(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	unit_id := c.Param("unit_id")
	// convert unit id to int
	unit_id_int, err := strconv.Atoi(unit_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid unit id"))
		return
	}
	employees, err := h.uc.GetByUnit(user_id.(string), unit_id_int)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get employee by unit", employees))
}

func (h *EmployeeHandler) Search(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	var payload struct {
		Input string `json:"input" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employees, err := h.uc.Search(user_id.(string), payload.Input)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Search employee", employees))
}

func (h *EmployeeHandler) Add(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	var payload domain.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employee, err := h.uc.Add(user_id.(string), &payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Add employee", employee))
}
