package service

import (
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RESTapi struct {
	Router *gin.Engine
}

func (r *RESTapi) Init(mode string) *gin.RouterGroup {
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Router = gin.New()
	r.Router.Use(gin.Logger())
	r.Router.Use(gin.Recovery())
	// r.Router = gin.Default()
	r.Router.Use(middleware.RateLimiter)
	r.Router.NoRoute(middleware.NoRouteExists)
	return r.Router.Group("/api/v1")
}
