package api

import (
	"net/http"

	"github.com/rick-you/uvid/portal"

	"github.com/labstack/echo/v4/middleware"
)

func bindDashStatic(server Server) {
	server.App.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(portal.DashFolder),
		Root:       "/packages/dash/dist",
		HTML5:      true,
	}))
}
