package api

import (
	"net/http"

	"luvsic3/uvid/portal"

	"github.com/labstack/echo/v4/middleware"
)

func bindDashStatic(server Server) {
	rg := server.App.Group("/dash")
	rg.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(portal.DashFolder),
		HTML5:      true,
	}))
}
