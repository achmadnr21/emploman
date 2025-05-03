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
	"fmt"
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
		// CRUD
		employee.GET("", empHandler.GetAll)
		employee.POST("", empHandler.Add)
		employee.PUT("/:nip", empHandler.UpdateEmployee)
		employee.POST("/uploadpp/:nip", empHandler.UploadPP)
		employee.GET("/:nip", empHandler.GetByNIP)
		employee.GET("/unit/:unit_id", empHandler.GetByUnit)
		employee.GET("/search", empHandler.Search)
		employee.PUT("/uprole/:nip", empHandler.Promote)
		employee.PUT("/downrole/:nip", empHandler.Demote)

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

// upload profile picture
// alurnya adalah akan menerima file jpeg/png dari client
// dan akan disimpan di s3 bucket

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

func (h *EmployeeHandler) GetMe(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	employee, err := h.uc.GetMe(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Get employee by NIP", employee))
}
func (h *EmployeeHandler) UpdateMe(c *gin.Context) {
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

func (h *EmployeeHandler) UploadPPMe(c *gin.Context) {
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

func (h *EmployeeHandler) Demote(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	var payload domain.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	employee, err := h.uc.Demote(user_id.(string), nip, payload.RoleID)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Demote employee", employee))
}
