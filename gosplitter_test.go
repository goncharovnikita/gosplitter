package gosplitter_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/goncharovnikita/gosplitter"
	"github.com/stretchr/testify/assert"
)

// Test ContextCaller
func caller() string {
	return gosplitter.CallerContext()
}
func TestContextCaller(t *testing.T) {
	result := caller()

	assert.Equal(t, "gosplitter_test.TestContextCaller", result)
}

func TestGetAbsoluteURL(t *testing.T) {
	gosplitter.RegisteredPatterns["testGetAbsolutePath"] = &gosplitter.RouterPoint{
		URL: "/test",
	}

	var result string
	var err error
	result, err = gosplitter.GetAbsoluteURL("/get", "testGetAbsolutePath")
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/get", result)
	result, err = gosplitter.GetAbsoluteURL("/test", "new")
	assert.Equal(t, "", result)
	assert.Equal(t, "pattern /test not registered", err.Error())
}

// Test RegisterHandler
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func handlerFunc(w http.ResponseWriter, r *http.Request) {}
func TestRegisterHandler(t *testing.T) {
	var (
		h      handler
		hf     = handlerFunc
		caller = "caller"
		mux    = http.NewServeMux()
		point  gosplitter.RouterPoint
	)
	gosplitter.RegisterHandler("/test1", caller, h, mux)

	assert.Equal(t, h, gosplitter.RegisteredPatterns[caller].Handler)

	gosplitter.RegisteredPatterns[caller] = nil

	gosplitter.RegisterHandler("/test2", caller, hf, mux)

	// t.Logf("%+v\n", gosplitter.RegisteredPatterns)

	point = *gosplitter.RegisteredPatterns[caller]
	assert.Equal(t, fmt.Sprintf("%v", hf), fmt.Sprintf("%v", point.HandlerFunc))

	gosplitter.RegisteredPatterns[caller] = &gosplitter.RouterPoint{
		URL: "/test",
	}

	gosplitter.RegisterHandler("/nested", caller, hf, mux)
	point = *gosplitter.RegisteredPatterns[caller]
	assert.Equal(t, "/test/nested", point.URL)

}

// Test registerRouterPoint
type rPoint struct{}

func TestRegisterRouterPoint(t *testing.T) {
	var err error
	var p rPoint
	gosplitter.RegisterRouterPoint(gosplitter.CallerContext(), "/test", p)

	assert.Equal(t, nil, err)
	// t.Logf("%v\n", gosplitter.RegisteredPatterns)
	assert.Equal(t, "/test", gosplitter.RegisteredPatterns["rPoint"].URL)
}
