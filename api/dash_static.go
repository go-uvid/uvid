package api

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	dash_embed "github.com/rick-you/uvid/packages"
)

func bindDashStatic(server Server) {
	server.App.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(dash_embed.DashFolder),
		Root:       "/packages/dash/dist",
		HTML5:      true,
	}))
}
