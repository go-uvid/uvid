package api

import (
	"luvsic3/uvid/daos"
	"luvsic3/uvid/dtos"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	App *echo.Echo
	Dao *daos.Dao
}

func New(dsn string) Server {
	server := Server{
		App: echo.New(),
		Dao: daos.New(dsn),
	}

	server.App.Validator = &dtos.CustomValidator{Validator: validator.New()}

	server.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	server.App.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(50)))
	server.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Minute,
	}))
	server.App.Use(middleware.Recover())
	server.App.Use(middleware.Secure())
	server.App.Use(middleware.Logger())

	bindSpanApi(server)
	return server
}

func (server Server) Start() {
	PORT := os.Getenv("PORT")
	server.App.Logger.Fatal(server.App.Start(":" + PORT))
}
