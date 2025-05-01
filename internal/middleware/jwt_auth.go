package middleware

import (
	"net/http"
	"strings"

	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("Token tidak valid!"))
		c.Abort()
		return
	}

	tokenString := bearerToken[1]
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("Token tidak valid!"))
		c.Abort()
		return
	}

	// Memvalidasi token
	claims, err := utils.ParseAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("Token tidak valid!"))
		c.Abort()
		return
	}

	// Menyimpan user_id di context Gin
	c.Set("user_id", claims.UserId)
	c.Set("issuer", claims.Issuer)
	c.Next()
}
