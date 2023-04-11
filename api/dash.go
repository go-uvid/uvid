package api

import (
	"luvsic3/uvid/tools"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func bindDashApi(server Server) {
	api := &dashApi{server: server}
	rg := server.App.Group("/dash")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: []byte(Configs["jwt_secret"]),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/dash/user/login"
		},
	}
	rg.Use(echojwt.WithConfig(config))

	rg.POST("/user/login", api.loginUser)
	rg.POST("/user/password", api.updateUserPassword)
	rg.GET("/pageview", api.pageview)
}

type dashApi struct {
	server Server
}

func (api *dashApi) loginUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := api.server.Dao.GetUserByName(username)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Throws unauthorized error
	if err := tools.ComparePassword(user.Password, password); err != nil {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(Configs["jwt_secret"]))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (api *dashApi) updateUserPassword(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (api *dashApi) pageview(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
