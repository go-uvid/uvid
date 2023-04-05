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
	rg.POST("/http", api.createHTTP, NewEnsureSessionMiddleware(api))
	rg.POST("/event", api.createEvent, NewEnsureSessionMiddleware(api))
	rg.POST("/performance", api.createPerformance, NewEnsureSessionMiddleware(api))
	rg.POST("/pageview", api.createPageView, NewEnsureSessionMiddleware(api))
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
	session, err := api.server.Dao.CreateSession(sessionDTO)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sessionUUID := session.UUID
	cookie := http.Cookie{
		Name:  SessionKey,
		Value: sessionUUID.String(),
	}
	c.SetCookie(&cookie)
	return c.NoContent(http.StatusNoContent)
}

func (api *spanApi) createError(c echo.Context) error {
	dto := &dtos.ErrorDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.server.Dao.CreateJSError(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createHTTP(c echo.Context) error {
	dto := &dtos.HTTPDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.server.Dao.CreateHTTP(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createEvent(c echo.Context) error {
	dto := &dtos.EventDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.server.Dao.CreateEvent(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createPerformance(c echo.Context) error {
	dto := &dtos.PerformanceDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.server.Dao.CreatePerformance(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createPageView(c echo.Context) error {
	dto := &dtos.PageViewDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.server.Dao.CreatePageView(session, dto)
	return handleDaoAndResponse(c, err)
}

func handleDaoAndResponse(c echo.Context, err error) error {
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
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
