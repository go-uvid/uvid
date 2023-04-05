package api_test

import (
	"luvsic3/uvid/api"
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tests"
	"luvsic3/uvid/tools"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	API_SESSION = "/span/session"
	API_ERROR   = "/span/error"
	API_HTTP    = "/span/http"
	API_EVENT   = "/span/event"
	API_PERF    = "/span/performance"
	API_PV      = "/span/pageview"
)

var requestHeader = map[string]string{
	"User-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
}

func getSessionApi(t *testing.T, sessionUUID *uuid.UUID, cookie *http.Cookie) tests.ApiScenario {
	var sessionApiScenarios = tests.ApiScenario{
		Name:           "create session",
		Method:         http.MethodPost,
		Url:            API_SESSION,
		ExpectedStatus: http.StatusNoContent,
		RequestHeaders: requestHeader,
		Body: strings.NewReader(`{
				"url": "www.google.com",
				"screen": "1080*700",
				"referrer": "www.google.com",
				"language": "en"
			  }`),
		AfterRequest: func(res *http.Response, server *api.Server) {
			cookies := res.Cookies()
			*cookie = *cookies[0]
			_uuid, err := uuid.Parse(cookie.Value)
			*sessionUUID = _uuid
			assert.NoError(t, err)
		},
	}

	return sessionApiScenarios
}

func TestCreateError(t *testing.T) {
	var sessionUUID uuid.UUID
	var cookie http.Cookie
	var sessionApiScenarios = getSessionApi(t, &sessionUUID, &cookie)

	var BeforeRequest = func(req *http.Request, server *api.Server) {
		assert.NotNil(t, cookie)
		assert.NotEqual(t, sessionUUID, uuid.Nil)
		req.AddCookie(&cookie)
	}

	scenarios := []tests.ApiScenario{
		sessionApiScenarios,
		{
			Name:           "create error",
			Method:         http.MethodPost,
			Url:            API_ERROR,
			Body:           strings.NewReader(tools.StructToJSONString(dtos.ErrorDTO{Name: "ErrorName", Message: "ErrorMessage", Stack: "ErrorStack"})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB.Find(&errors)
				assert.Len(t, errors, 1)
				assert.Equal(t, errors[0].SessionUUID, sessionUUID)
			},
		},
		{
			Name:           "create error 2",
			Method:         http.MethodPost,
			Url:            API_ERROR,
			Body:           strings.NewReader(tools.StructToJSONString(dtos.ErrorDTO{Name: "ErrorName2", Message: "ErrorMessage2", Stack: "ErrorStack2"})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB.Find(&errors)
				assert.Len(t, errors, 2)
				assert.Equal(t, errors[0].SessionUUID, sessionUUID)
				assert.Equal(t, errors[1].SessionUUID, sessionUUID)
			},
		},
		{
			Name:   "create http",
			Method: http.MethodPost,
			Url:    API_HTTP,
			Body: strings.NewReader(tools.StructToJSONString(dtos.HTTPDTO{
				Resource: randomReferrer(),
				Method:   randomHttpMethod(),
				Headers:  "Content-Type: application/json",
				Status:   http.StatusInternalServerError,
				Data:     "",
				Response: "",
			})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.HTTP{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionUUID, sessionUUID)
			},
		},
		{
			Name:   "create event",
			Method: http.MethodPost,
			Url:    API_EVENT,
			Body: strings.NewReader(tools.StructToJSONString(dtos.EventDTO{
				Name:  randomEventName(),
				Value: "",
			})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.Event{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionUUID, sessionUUID)
			},
		},
		{
			Name:   "create performance",
			Method: http.MethodPost,
			Url:    API_PERF,
			Body: strings.NewReader(tools.StructToJSONString(dtos.PerformanceDTO{
				Name:  randomPerfName(),
				Value: randomPerfValue(),
				URL:   randomReferrer(),
			})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.Performance{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionUUID, sessionUUID)
			},
		},
		{
			Name:   "create pageview",
			Method: http.MethodPost,
			Url:    API_PV,
			Body: strings.NewReader(tools.StructToJSONString(dtos.PerformanceDTO{
				URL: randomReferrer(),
			})),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.PageView{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionUUID, sessionUUID)
			},
		},
	}

	testCase := tests.NewTestCase(*t, scenarios, nil, nil)
	testCase.Test()
}

func TestSessionMiddleware(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:           "create event without session",
			Method:         http.MethodPost,
			Url:            API_EVENT,
			Body:           strings.NewReader(`{"name": "register", "value": "new user"}`),
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	testCase := tests.NewTestCase(*t, scenarios, nil, nil)
	testCase.Test()
}

// randomPerfName generates a random name for a PerformanceMetric
func randomPerfName() string {
	names := []string{models.LCP, models.CLS, models.FID}
	return names[rand.Intn(len(names))]
}

// randomPerfValue generates a random value for a PerformanceMetric
func randomPerfValue() float64 {
	return rand.Float64() * 10
}

// randomReferrer generates a random URL for a PerformanceMetric
func randomReferrer() string {
	urls := []string{"https://example.com", "https://google.com", "https://github.com", "https://stackoverflow.com", "https://wikipedia.org"}
	return urls[rand.Intn(len(urls))]
}

func randomHttpMethod() string {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect, http.MethodTrace}
	return methods[rand.Intn(len(methods))]
}

func randomEventName() string {
	names := []string{"register", "login", "logout", "click", "view", "add", "remove", "update"}
	return names[rand.Intn(len(names))]
}
