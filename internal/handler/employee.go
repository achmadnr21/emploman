package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/middleware"
	usecase "github.com/achmadnr21/emploman/internal/usecase/employee"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	uc *usecase.EmployeeUsecase
}

func NewEmployeeHandler(apiV *gin.RouterGroup, uc *usecase.EmployeeUsecase) {
	EmployeeHandler := &EmployeeHandler{
		uc: uc,
	}

	employee := apiV.Group("/employee")
	employee.Use(middleware.JWTAuthMiddleware)
	{
		// Basic CRUD
		employee.GET("", EmployeeHandler.GetAll)              // GET /employees
		employee.POST("", EmployeeHandler.Add)                // POST /employees
		employee.GET("/:nip", EmployeeHandler.GetByNIP)       // GET /employees/:nip
		employee.PUT("/:nip", EmployeeHandler.UpdateEmployee) // PUT /employees/:nip

		// Upload profile picture
		employee.POST("/:nip/profile-picture", EmployeeHandler.UploadPP) // POST /employees/:nip/profile-picture

		// Filtering
		employee.GET("/unit/:unit_id", EmployeeHandler.GetByUnit) // GET /employees/unit/:unit_id
		employee.GET("/search", EmployeeHandler.Search)           // GET /employees/search

		// Promotion
		employee.PUT("/:nip/promote", EmployeeHandler.Promote) // PUT /employees/:nip/promote

	}
}

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
	// contoh endpoint: /employee/search?query=rudy traspac
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Query is required"))
		return
	}
	employees, err := h.uc.Search(user_id.(string), query)
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

func (h *EmployeeHandler) UploadPP(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid file"))
		return
	}
	url, err := h.uc.UploadPP(user_id.(string), nip, file)
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

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	var payload domain.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employee, err := h.uc.UpdateEmployee(user_id.(string), nip, &payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Update employee", employee))
}

func (h *EmployeeHandler) Promote(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	var payload domain.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employee, err := h.uc.Promote(user_id.(string), nip, payload.RoleID)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Promote employee", employee))
}
