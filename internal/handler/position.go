package handler

import (
	"net/http"
	"strconv"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	uc *usecase.PositionUsecase
}

func NewPositionHandler(uc *usecase.PositionUsecase) *PositionHandler {
	return &PositionHandler{
		uc: uc,
	}
}

/*
position := apiV.Group("/position")

	position.Use(middleware.JWTAuthMiddleware)
	{
		position.GET("", positionHandler.GetAllPosition) // GET /positions
		position.POST("", positionHandler.AddPosition)   // POST /positions
		position.GET("/:id", positionHandler.GetPositionByID)
		position.PUT("/:id", positionHandler.UpdatePosition)
		position.DELETE("/:id", positionHandler.DeletePosition)
		position.GET("/search", positionHandler.SearchPosition) // GET /positions/search
	}
*/

func (h *PositionHandler) AddPosition(c *gin.Context) {
	user_id := c.MustGet("user_id").(string)
	var position domain.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError("Invalid input"))
		return
	}
	newposition, err := h.uc.AddPosition(user_id, &position)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, newposition)
}
func (h *PositionHandler) UpdatePosition(c *gin.Context) {
	user_id := c.MustGet("user_id").(string)
	id := c.Param("id")
	// convert id to int
	idInt, err := strconv.Atoi(id)
	var position domain.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError("Invalid input"))
		return
	}
	position.ID = idInt
	updatedPosition, err := h.uc.UpdatePosition(user_id, &position)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, updatedPosition)
}
func (h *PositionHandler) DeletePosition(c *gin.Context) {
	user_id := c.MustGet("user_id").(string)
	id := c.Param("id")
	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid id"))
		return
	}
	err = h.uc.DeletePosition(user_id, idInt)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Position deleted", nil))
}
func (h *PositionHandler) GetAllPosition(c *gin.Context) {
	positions, err := h.uc.GetAllPosition()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, positions)
}
func (h *PositionHandler) GetPositionByID(c *gin.Context) {
	id := c.Param("id")
	// convert id to int
	idInt, err := strconv.Atoi(id)
	position, err := h.uc.GetPositionByID(idInt)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, position)
}

func (h *PositionHandler) SearchPosition(c *gin.Context) {
	var payload struct {
		Query string `form:"query" json:"query"`
	}
	if err := c.ShouldBindQuery(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid input"))
		return
	}
	positions, err := h.uc.SearchPosition(payload.Query)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, positions)
}
