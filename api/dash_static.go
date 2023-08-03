package api

import (
	"net/http"

	dash_embed "github.com/go-uvid/uvid/js"
	"github.com/labstack/echo/v4/middleware"
)

func bindDashStatic(server Server) {
	server.App.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(dash_embed.DashFolder),
		Root:       "/packages/dash/dist",
		HTML5:      true,
	}))
}
