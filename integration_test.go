package gosplitter_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goncharovnikita/gosplitter"
)

type UsersRouterPoint struct{}

func (u UsersRouterPoint) Greeter(mux *http.ServeMux) {
	var h GreeterHandler
	gosplitter.Match("/hello", mux, h)
}

func (u UsersRouterPoint) Start(mux *http.ServeMux) {
	u.Greeter(mux)
}

type GreeterHandler struct{}

func (c GreeterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

type APIV1Handler struct {
	mux *http.ServeMux
}

func (a *APIV1Handler) Start() {
	var colorsHandler = ColorsHandler{
		mux: a.mux,
	}
	gosplitter.Match("/ping", a.mux, a.HandlePing())
	gosplitter.Match("/hru", a.mux, a.HandleMood())
	gosplitter.Match("/colors", a.mux, colorsHandler)
	colorsHandler.Start()
}

func (a *APIV1Handler) HandlePing() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func (a *APIV1Handler) HandleMood() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nice"))
	}
}

type ColorsHandler struct {
	mux *http.ServeMux
}

func (c *ColorsHandler) Start() {
	gosplitter.Match("/black", c.mux, c.HandleBlack())
}

func (c *ColorsHandler) HandleBlack() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("#000000"))
	}
}

func TestIntegration(t *testing.T) {
	var u UsersRouterPoint
	var err error
	var mux = http.NewServeMux()
	var response *http.Response
	var body []byte
	var c http.Client
	var apiV1 = APIV1Handler{
		mux: mux,
	}
	gosplitter.Match("/api", mux, u)
	gosplitter.Match("/api/v1", mux, apiV1)
	u.Start(mux)
	apiV1.Start()
	// fmt.Printf("%+v\n", gosplitter.RegisteredPatterns)
	go http.ListenAndServe(":8080", mux)

	assert.Equal(t, nil, err)

	// /api/hello
	response, err = c.Get("http://localhost:8080/api/hello")
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello!", string(body))

	// /api/v1/ping
	response, err = c.Get("http://localhost:8080/api/v1/ping")
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "pong", string(body))

	// /api/v1/hru
	response, err = c.Get("http://localhost:8080/api/v1/hru")
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "nice", string(body))

	// /api/v1/colors/black
	response, err = c.Get("http://localhost:8080/api/v1/colors/black")
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "#000000", string(body))
}
