package api_test

import (
	"luvsic3/uvid/api"
	"luvsic3/uvid/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateError(t *testing.T) {
	server := NewTestServer()
	rec := server.requestSession()
	cookies := rec.Result().Cookies()
	cookie := cookies[0]
	_uuid, err := uuid.Parse(cookie.Value)
	assert.NoError(t, err)

	// first error
	body1 := `{"name": "ErrorName", "message": "ErrorMessage", "stack": "ErrorStack"}`
	rec = server.request(body1, cookie)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	errors := []models.JSError{}
	server.Dao.DB().Find(&errors)
	assert.Len(t, errors, 1)
	assert.Equal(t, errors[0].SessionUUID, _uuid)

	// second error
	body2 := `{"name": "ErrorName2", "message": "ErrorMessage2", "stack": "ErrorStack2"}`
	rec = server.request(body2, cookie)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	errors = []models.JSError{}
	server.Dao.DB().Find(&errors)
	assert.Len(t, errors, 2)
	assert.Equal(t, errors[0].SessionUUID, _uuid)
	assert.Equal(t, errors[1].SessionUUID, _uuid)
}

type TestServer struct {
	api.Server
}

func NewTestServer() TestServer {
	server := api.New(":memory:")
	return TestServer{
		server,
	}
}

func (server TestServer) requestSession() *httptest.ResponseRecorder {
	body := `{
		"url": "www.google.com",
		"screen": "1080*700",
		"referrer": "www.google.com",
		"language": "en"
	  }`
	req, rec := makePostRequest("/span/session", body)
	server.App.ServeHTTP(rec, req)
	return rec
}

func (server TestServer) request(body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req, rec := makePostRequest("/span/error", body)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	server.App.ServeHTTP(rec, req)
	return rec
}

func makePostRequest(path string, body string) (*http.Request, *httptest.ResponseRecorder) {
	reader := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, path, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // Set the Content-Type header
	req.Header.Set("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	rec := httptest.NewRecorder()
	return req, rec
}
