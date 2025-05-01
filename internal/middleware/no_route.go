package middleware

import (
	"net/http"

	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

func NoRouteExists(c *gin.Context) {
	c.JSON(http.StatusNotFound, utils.ResponseError("No Route Exists"))
}
