package handler

import (
	"net/http"
	"strconv"

	usecase "github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type PrintHandler struct {
	uc *usecase.PrintUsecase
}

func NewPrintHandler(uc *usecase.PrintUsecase) *PrintHandler {
	return &PrintHandler{
		uc: uc,
	}
}

func (h *PrintHandler) PrintAll(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	employees, err := h.uc.PrintAll(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Print all employees", employees))
}
func (h *PrintHandler) PrintByUnitID(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	unit_id := c.Param("unit_id")
	// convert unit id to int
	unit_id_int, err := strconv.Atoi(unit_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid unit id"))
		return
	}
	employees, err := h.uc.PrintByUnitID(user_id.(string), unit_id_int)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Print employee by unit", employees))
}
func (h *PrintHandler) PrintByNIP(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	nip := c.Param("nip")
	employee, err := h.uc.PrintByNIP(user_id.(string), nip)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Print employee by NIP", employee))
}
