package gosplitter_test

import (
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

	assert.Equal(t, "github.com/goncharovnikita/gosplitter_test.TestContextCaller", result)
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

type handlerFunc struct{}

func (h handlerFunc) Handle(w http.ResponseWriter, r *http.Request) {}

func TestRegisterHandler(t *testing.T) {
	var h handler
	var hf handlerFunc
	var err error
	var mux = http.NewServeMux()
	var testWithChildrenOne = "twc1"
	var testWithChildrenTwo = "twc2"
	/**
	* Test register Handler without children
	 */
	err = gosplitter.RegisterHandler("/test", "testRegisterHandler", h, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, h, gosplitter.RegisteredPatterns["testRegisterHandler"].Handler)

	// Test with children case one
	gosplitter.RegisteredPatterns[testWithChildrenOne] = &gosplitter.RouterPoint{
		URL: "/test/child1",
	}
	err = gosplitter.RegisterHandler("/handler", testWithChildrenOne, h, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, h, gosplitter.RegisteredPatterns[testWithChildrenOne].Children["/handler"].Handler)

	err = gosplitter.RegisterHandler("/handler/func", testWithChildrenOne, hf, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, hf, gosplitter.RegisteredPatterns[testWithChildrenOne].Children["/handler/func"].HandlerFunc)

	// Test with children case two
	gosplitter.RegisteredPatterns[testWithChildrenTwo] = &gosplitter.RouterPoint{
		URL: "/test/child2",
	}

	err = gosplitter.RegisterHandler("/handler/func", testWithChildrenTwo, hf, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, hf, gosplitter.RegisteredPatterns[testWithChildrenTwo].Children["/handler/func"].HandlerFunc)

	err = gosplitter.RegisterHandler("/handler", testWithChildrenTwo, h, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, h, gosplitter.RegisteredPatterns[testWithChildrenTwo].Children["/handler"].Handler)

	err = gosplitter.RegisterHandler("/test", "testRegisterHandler", h, mux)
	assert.Equal(t, "pattern /test already registered", err.Error())

	err = gosplitter.RegisterHandler("/test2", "testRegisterHandlerFunc", hf, mux)

	assert.Equal(t, nil, err)
	assert.Equal(t, hf, gosplitter.RegisteredPatterns["testRegisterHandlerFunc"].HandlerFunc)
}

// Test registerRouterPoint
func TestRegisterRouterPoint(t *testing.T) {
	var err error
	err = gosplitter.RegisterRouterPoint("/test", "routerPoint")

	assert.Equal(t, nil, err)
	assert.Equal(t, "/test", gosplitter.RegisteredPatterns["routerPoint"].Children["/test"].URL)
}
