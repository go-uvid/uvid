package api

import (
	"net/http"

	"github.com/rick-you/uvid/dtos"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"zgo.at/isbot"
)

func bindSpanApi(server Server) {
	api := &spanApi{server}
	rg := server.App.Group("/span", checkBotMiddleware)
	rg.POST("/session", api.createSession)
	rg.POST("/error", api.createError)
	rg.POST("/http", api.createHTTP)
	rg.POST("/event", api.createEvent)
	rg.POST("/performance", api.createPerformance)
	rg.POST("/pageview", api.createPageView)
}

type spanApi struct {
	Server
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
	session, err := api.Dao.CreateSession(sessionDTO)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sessionUUID := session.UUID
	return c.String(http.StatusOK, sessionUUID.String())
}

func (api *spanApi) createError(c echo.Context) error {
	dto := &dtos.ErrorDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.Dao.CreateJSError(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createHTTP(c echo.Context) error {
	dto := &dtos.HTTPDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.Dao.CreateHTTP(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createEvent(c echo.Context) error {
	dto := &dtos.EventDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.Dao.CreateEvent(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createPerformance(c echo.Context) error {
	dto := &dtos.PerformanceDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.Dao.CreatePerformance(session, dto)
	return handleDaoAndResponse(c, err)
}

func (api *spanApi) createPageView(c echo.Context) error {
	dto := &dtos.PageViewDTO{}
	if err := dtos.BindAndValidateDTO(c, dto); err != nil {
		return err
	}
	session := GetSessionUUID(c)
	_, err := api.Dao.CreatePageView(session, dto)
	return handleDaoAndResponse(c, err)
}

func handleDaoAndResponse(c echo.Context, err error) error {
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusNoContent)
}

const SessionHeaderKey = "X-UVID-Session"

func GetSessionUUID(c echo.Context) uuid.UUID {
	return uuid.MustParse(c.Request().Header.Get(SessionHeaderKey))
}

func checkBotMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		result := isbot.Bot(c.Request())
		if isbot.Is(result) {
			return echo.ErrForbidden
		}
		return next(c)
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
