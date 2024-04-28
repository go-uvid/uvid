package api_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/go-uvid/uvid/api"
	"github.com/go-uvid/uvid/dtos"
	"github.com/go-uvid/uvid/models"
	"github.com/go-uvid/uvid/tests"
	"github.com/go-uvid/uvid/tools"

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

func getSessionApi(t *testing.T, sessionUUID *uint) tests.ApiScenario {
	var sessionApiScenarios = tests.ApiScenario{
		Name:           "create session",
		Method:         http.MethodPost,
		Url:            API_SESSION,
		ExpectedStatus: http.StatusOK,
		RequestHeaders: requestHeader,
		Body: strings.NewReader(`{
				"url": "www.google.com",
				"screen": "1080*700",
				"referrer": "www.google.com",
				"language": "en"
			  }`),
		AfterRequest: func(res *http.Response, server *api.Server) {
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			_id, err := strconv.Atoi(string(body))
			assert.NoError(t, err)
			*sessionUUID = uint(_id)
		},
	}

	return sessionApiScenarios
}

func TestCreateSpan(t *testing.T) {
	var sessionID uint
	var sessionApiScenarios = getSessionApi(t, &sessionID)

	var BeforeRequest = func(req *http.Request, server *api.Server) {
		assert.NotEqual(t, sessionID, uuid.Nil)
		req.Header.Set(api.SessionHeaderKey, fmt.Sprint(sessionID))
	}

	var createSpanApiScenario = func(scenario tests.ApiScenario) tests.ApiScenario {
		return tests.ApiScenario{
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest:  BeforeRequest,
			RequestHeaders: requestHeader,
			Body:           scenario.Body,
			Url:            scenario.Url,
			AfterRequest:   scenario.AfterRequest,
		}
	}

	scenarios := []tests.ApiScenario{
		sessionApiScenarios,
		createSpanApiScenario(tests.ApiScenario{
			Name: "create error",
			Url:  API_ERROR,
			Body: strings.NewReader(tools.StructToJSONString(dtos.ErrorDTO{Name: "ErrorName", Message: "ErrorMessage", Stack: "ErrorStack"})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB.Find(&errors)
				assert.Len(t, errors, 1)
				assert.Equal(t, errors[0].SessionID, sessionID)
			},
		}),
		createSpanApiScenario(tests.ApiScenario{
			Name: "create error 2",
			Url:  API_ERROR,
			Body: strings.NewReader(tools.StructToJSONString(dtos.ErrorDTO{Name: "ErrorName2", Message: "ErrorMessage2", Stack: "ErrorStack2"})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB.Find(&errors)
				assert.Len(t, errors, 2)
				assert.Equal(t, errors[0].SessionID, sessionID)
				assert.Equal(t, errors[1].SessionID, sessionID)
			},
		}),
		createSpanApiScenario(tests.ApiScenario{
			Name: "create http",
			Url:  API_HTTP,
			Body: strings.NewReader(tools.StructToJSONString(dtos.HTTPDTO{
				Resource: randomReferrer(),
				Method:   randomHttpMethod(),
				Headers:  "Content-Type: application/json",
				Status:   http.StatusInternalServerError,
				Body:     "",
				Response: "",
			})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.HTTP{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionID, sessionID)
			},
		}),
		createSpanApiScenario(tests.ApiScenario{
			Name: "create event",
			Url:  API_EVENT,
			Body: strings.NewReader(tools.StructToJSONString(dtos.EventDTO{
				Action: randomEventAction(),
				Value:  "",
			})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.Event{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionID, sessionID)
			},
		}),
		createSpanApiScenario(tests.ApiScenario{
			Name: "create performance",
			Url:  API_PERF,
			Body: strings.NewReader(tools.StructToJSONString(dtos.PerformanceDTO{
				Name:  randomPerfName(),
				Value: randomPerfValue(),
				URL:   randomReferrer(),
			})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.Performance{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionID, sessionID)
			},
		}),
		createSpanApiScenario(tests.ApiScenario{
			Name: "create pageview",
			Url:  API_PV,
			Body: strings.NewReader(tools.StructToJSONString(dtos.PerformanceDTO{
				URL: randomReferrer(),
			})),
			AfterRequest: func(res *http.Response, server *api.Server) {
				https := []models.PageView{}
				server.Dao.DB.Find(&https)
				assert.Len(t, https, 1)
				assert.Equal(t, https[0].SessionID, sessionID)
			},
		}),
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
			Body:           strings.NewReader(`{"action": "register", "value": "new user"}`),
			ExpectedStatus: http.StatusInternalServerError,
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

func randomEventAction() string {
	names := []string{"register", "login", "logout", "click", "view", "add", "remove", "update"}
	return names[rand.Intn(len(names))]
}
