package api

import (
	"net/http"
	"time"

	"github.com/rick-you/uvid/dtos"
	"github.com/rick-you/uvid/tools"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func bindDashApi(server Server) {
	api := &dashApi{server}
	rg := server.App.Group("/dash")
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
	rg.POST("/user/password", api.changeUserPassword)

	rg.GET("/pvs", api.pageviews)
	rg.GET("/pvs/interval", api.pageviewInterval)
	rg.GET("/pvs/count", api.pageviewCount)

	rg.GET("/uvs/interval", api.uniqueVisitorInterval)
	rg.GET("/uvs/count", api.uniqueVisitorCount)

	rg.GET("/errors", api.errors)
	rg.GET("/errors/interval", api.errorInterval)
	rg.GET("/errors/count", api.errorCount)

	rg.GET("/https/errors", api.httpErrorInterval)
	rg.GET("/https/errors/interval", api.httpErrorInterval)
	rg.GET("/https/errors/count", api.httpErrorCount)

	rg.GET("/sessions", api.sessions)
	rg.GET("/performances", api.avgPerformance)
	rg.GET("/events/group", api.eventGroup)
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
		return echo.ErrBadRequest
	}
	// Throws unauthorized error
	if err := tools.ComparePassword(user.Password, body.Password); err != nil {
		return echo.ErrBadRequest
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

func (api *dashApi) changeUserPassword(c echo.Context) error {
	body := &dtos.ChangePasswordDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	if err := api.Dao.ChangeUserPassword(name, body.CurrentPassword, body.NewPassword); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (api *dashApi) pageviews(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	pageviews, err := api.Dao.FindPageViews(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pageviews)
}

func (api *dashApi) pageviewInterval(c echo.Context) error {
	body := &dtos.TimeIntervalSpanDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindPageViewInterval(api.Dao.SpanFilter(body.Start, body.End), body.Unit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) pageviewCount(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	count, err := api.Dao.FindPageViewCount(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, count)
}

func (api *dashApi) uniqueVisitorInterval(c echo.Context) error {
	body := &dtos.TimeIntervalSpanDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindUniqueVisitorInterval(api.Dao.SpanFilter(body.Start, body.End), body.Unit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) uniqueVisitorCount(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	count, err := api.Dao.FindUniqueVisitorCount(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, count)
}

func (api *dashApi) avgPerformance(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindAveragePerformanceInterval(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) sessions(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	sessions, err := api.Dao.FindSessions(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, sessions)
}

func (api *dashApi) eventGroup(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindEventInterval(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) errors(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	errors, err := api.Dao.FindJSErrors(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, errors)
}

func (api *dashApi) errorInterval(c echo.Context) error {
	body := &dtos.TimeIntervalSpanDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindJSErrorInterval(api.Dao.SpanFilter(body.Start, body.End), body.Unit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) errorCount(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindJSErrorCount(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) httpErrors(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	errors, err := api.Dao.FindHTTPErrors(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, errors)
}

func (api *dashApi) httpErrorInterval(c echo.Context) error {
	body := &dtos.TimeIntervalSpanDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindHTTPErrorInterval(api.Dao.SpanFilter(body.Start, body.End), body.Unit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}

func (api *dashApi) httpErrorCount(c echo.Context) error {
	body := &dtos.SpanFilterDTO{}
	if err := dtos.BindAndValidateDTO(c, body); err != nil {
		return err
	}

	interval, err := api.Dao.FindHTTPErrorCount(api.Dao.SpanFilter(body.Start, body.End))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, interval)
}
