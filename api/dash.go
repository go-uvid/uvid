package api

import (
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/tools"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func bindDashApi(server Server) {
	api := &dashApi{server}
	rg := server.App.Group("/dash")
	rg.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
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
	Server
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (api *dashApi) loginUser(c echo.Context) error {
	body := &dtos.LoginDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	user, err := api.Dao.GetUserByName(body.Name)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Throws unauthorized error
	if err := tools.ComparePassword(user.Password, body.Password); err != nil {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		body.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
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
	body := &dtos.UpdatePasswordDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	api.Dao.UpdateUserPassword(name, body.Password)
	return c.NoContent(http.StatusNoContent)
}

func (api *dashApi) pageview(c echo.Context) error {
	body := &dtos.TimeRangeDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindPageViewInterval(api.Dao.TimeRange(body.Start, body.End), body.Unit == dtos.UnitHour)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}
