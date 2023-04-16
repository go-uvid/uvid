package api

import (
	"luvsic3/uvid/daos"
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/tools"
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

var Configs = map[string]string{}

func New(dsn string) Server {
	server := Server{
		App: echo.New(),
		Dao: daos.New(dsn),
	}
	// daos.Seed(dsn)
	err := server.Dao.InitializeDB()
	if err != nil {
		panic(err)
	}
	Configs, err = server.Dao.GetAllConfigs()
	if err != nil {
		panic(err)
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
	bindDashApi(server)
	bindDashStatic(server)
	return server
}

func (server Server) Start() {
	envPort := os.Getenv("PORT")
	port := tools.Ternary(len(envPort) == 0, "3000", envPort)
	server.App.Logger.Fatal(server.App.Start(":" + port))
}
