package api

import (
	"luvsic3/uvid/dtos"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func bindSpanApi(server Server) {
	api := &spanApi{server: server}
	rg := server.App.Group("/span")
	rg.POST("/session", api.createSession)
	rg.POST("/error", api.createError, NewEnsureSessionMiddleware(api))
}

type spanApi struct {
	server Server
}

func (api *spanApi) createSession(c echo.Context) error {
	sessionDTO := createSessionDTOFromContext(c)
	if sessionDTO == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid sessionDTO")
	}
	if err := c.Validate(sessionDTO); err != nil {
		c.Logger().Error(err)
		return err
	}
	sessionUUID := api.server.Dao.CreateSession(sessionDTO)
	cookie := http.Cookie{
		Name:  SessionKey,
		Value: sessionUUID.String(),
	}
	c.SetCookie(&cookie)
	return c.NoContent(http.StatusNoContent)
}

func (api *spanApi) createError(c echo.Context) error {
	errorDTO := dtos.ErrorDTO{}
	if err := c.Bind(&errorDTO); err != nil {
		c.Logger().Error(err)
		return err
	}
	if err := c.Validate(errorDTO); err != nil {
		c.Logger().Error(err)
		return err
	}
	session := GetSessionUUID(c)
	api.server.Dao.CreateJSError(session, &errorDTO)
	return c.NoContent(http.StatusNoContent)
}

const SessionKey = "uvid-session"

func GetSessionUUID(c echo.Context) uuid.UUID {
	return c.Get(SessionKey).(uuid.UUID)
}

// NewEnsureSessionMiddleware ensures that the session cookie exist in the request
// SDK should carry the session header in all requests to create Span
func NewEnsureSessionMiddleware(api *spanApi) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionCookie, err := c.Cookie(SessionKey)
			if err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusBadRequest, "No session")
			}

			sessionUUID, err := uuid.Parse(sessionCookie.Value)
			if err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid session")
			}

			c.Set(SessionKey, sessionUUID)
			return next(c)
		}
	}
}

func createSessionDTOFromContext(c echo.Context) *dtos.SessionDTO {
	baseSession := dtos.BaseSessionDTO{}
	if err := c.Bind(&baseSession); err != nil {
		c.Logger().Error(err)
		return nil
	}
	if err := c.Validate(baseSession); err != nil {
		c.Logger().Error(err)
		return nil
	}

	return &dtos.SessionDTO{
		UA:             c.Request().UserAgent(),
		IP:             c.RealIP(),
		BaseSessionDTO: baseSession,
	}
}
