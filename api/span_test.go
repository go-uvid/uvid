package api_test

import (
	"luvsic3/uvid/api"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tests"
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
)

var requestHeader = map[string]string{
	"User-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
}

func TestCreateError(t *testing.T) {
	var sessionUUID uuid.UUID
	var cookie *http.Cookie
	scenarios := []tests.ApiScenario{
		{
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
				cookie = cookies[0]
				_uuid, err := uuid.Parse(cookie.Value)
				sessionUUID = _uuid
				assert.NoError(t, err)
			},
		},
		{
			Name:           "create error",
			Method:         http.MethodPost,
			Url:            API_ERROR,
			Body:           strings.NewReader(`{"name": "ErrorName", "message": "ErrorMessage", "stack": "ErrorStack"}`),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest: func(req *http.Request, server *api.Server) {
				req.AddCookie(cookie)
			},
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB().Find(&errors)
				assert.Len(t, errors, 1)
				assert.Equal(t, errors[0].SessionUUID, sessionUUID)
			},
		},
		{
			Name:           "create error 2",
			Method:         http.MethodPost,
			Url:            API_ERROR,
			Body:           strings.NewReader(`{"name": "ErrorName2", "message": "ErrorMessage2", "stack": "ErrorStack2"}`),
			ExpectedStatus: http.StatusNoContent,
			BeforeRequest: func(req *http.Request, server *api.Server) {
				req.AddCookie(cookie)
			},
			AfterRequest: func(res *http.Response, server *api.Server) {
				errors := []models.JSError{}
				server.Dao.DB().Find(&errors)
				assert.Len(t, errors, 2)
				assert.Equal(t, errors[0].SessionUUID, sessionUUID)
				assert.Equal(t, errors[1].SessionUUID, sessionUUID)
			},
		},
	}

	testCase := tests.NewTestCase(*t, scenarios, nil, nil)
	testCase.Test()
}

func TestSessionMiddleware(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:           "create event",
			Method:         http.MethodPost,
			Url:            API_EVENT,
			Body:           strings.NewReader(`{"name": "register", "value": "new user"}`),
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	testCase := tests.NewTestCase(*t, scenarios, nil, nil)
	testCase.Test()
}
