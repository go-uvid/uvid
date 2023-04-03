package tests

import (
	"context"
	"fmt"
	"io"
	"luvsic3/uvid/api"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

// ApiScenario defines a single api request test case/scenario.
type ApiScenario struct {
	Name           string
	Method         string
	Url            string
	Body           io.Reader
	RequestHeaders map[string]string

	// Delay adds a delay before checking the expectations usually
	// to ensure that all fired non-awaited go routines have finished
	Delay time.Duration

	// expectations
	// ---
	ExpectedStatus  int
	ExpectedContent interface{}

	// test hooks
	// ---
	BeforeRequest func(req *http.Request, server *Server)
	AfterRequest  func(res *http.Response, server *Server)
}

type Server = api.Server

type ApiTestCase struct {
	t         testing.T
	scenarios []ApiScenario
	server    Server
	// test hooks
	// ---
	BeforeTestFunc func(t *testing.T, app *Server)
	AfterTestFunc  func(t *testing.T, app *Server)
}

func NewTestCase(t testing.T, scenarios []ApiScenario, beforeTestFunc, afterTestFunc func(t *testing.T, app *Server)) ApiTestCase {
	server := api.New(":memory:")
	return ApiTestCase{
		t,
		scenarios,
		server,
		beforeTestFunc,
		afterTestFunc,
	}
}

// Test executes the test case/scenario.
func (testCase *ApiTestCase) Test() {
	testServer := &testCase.server

	if testCase.BeforeTestFunc != nil {
		testCase.BeforeTestFunc(&testCase.t, testServer)
	}

	// add middleware to timeout long-running requests (eg. keep-alive routes)
	testServer.App.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancelFunc := context.WithTimeout(c.Request().Context(), 100*time.Millisecond)
			defer cancelFunc()
			c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	})

	for _, v := range testCase.scenarios {
		testCase.request(&v)
	}

	if testCase.AfterTestFunc != nil {
		testCase.AfterTestFunc(&testCase.t, testServer)
	}
}

func (testCase *ApiTestCase) request(apiScenario *ApiScenario) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(apiScenario.Method, apiScenario.Url, apiScenario.Body)

	if apiScenario.BeforeRequest != nil {
		apiScenario.BeforeRequest(req, &testCase.server)
	}

	// set default header
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// set scenario headers
	for k, v := range apiScenario.RequestHeaders {
		req.Header.Set(k, v)
	}

	// execute request
	testCase.server.App.ServeHTTP(recorder, req)

	res := recorder.Result()

	if apiScenario.AfterRequest != nil {
		apiScenario.AfterRequest(res, &testCase.server)
	}

	var prefix = apiScenario.Name
	if prefix == "" {
		prefix = fmt.Sprintf("%s:%s", apiScenario.Method, apiScenario.Url)
	}

	if res.StatusCode != apiScenario.ExpectedStatus {
		testCase.t.Errorf("[%s] Expected status code %d, got %d", prefix, apiScenario.ExpectedStatus, res.StatusCode)
	}

	if apiScenario.Delay > 0 {
		time.Sleep(apiScenario.Delay)
	}

	if apiScenario.ExpectedContent != nil {
		// TODO io.ReadCloser
		assertSameShape(&testCase.t, apiScenario.ExpectedContent, res.Body)
	}
}

func assertSameShape(t *testing.T, expected interface{}, actual interface{}) {
	expectedReflectType := reflect.TypeOf(expected)
	actualReflectType := reflect.TypeOf(actual)

	if expectedReflectType.Kind() != actualReflectType.Kind() {
		t.Errorf("expected kind %v but got %v", expectedReflectType.Kind(), actualReflectType.Kind())
	}

	expectedFields := getFieldNames(expected)
	actualFields := getFieldNames(actual)

	if !reflect.DeepEqual(expectedFields, actualFields) {
		t.Errorf("expected fields %v but got %v", expectedFields, actualFields)
	}
}

func getFieldNames(s interface{}) []string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embedded := getFieldNames(reflect.New(field.Type).Interface())
			fields = append(fields, embedded...)
		} else {
			fields = append(fields, field.Name)
		}
	}

	return fields
}
