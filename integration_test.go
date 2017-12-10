package gosplitter_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goncharovnikita/gosplitter"
)

type UsersRouterPoint struct{}
type CreateUsersHandler struct{}

func (c CreateUsersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
func (u UsersRouterPoint) ServeUsers(mux *http.ServeMux) {
	var h CreateUsersHandler
	gosplitter.Match("/create", mux, h)
}

func (u UsersRouterPoint) Activate(mux *http.ServeMux) {
	u.ServeUsers(mux)
}
func TestIntegration(t *testing.T) {
	var u UsersRouterPoint
	var err error
	var mux = http.NewServeMux()
	var response *http.Response
	var body []byte
	var c http.Client
	err = gosplitter.Match("/api", mux, u)
	u.Activate(mux)
	t.Logf("%v", gosplitter.RegisteredPatterns)
	go http.ListenAndServe(":8080", mux)

	assert.Equal(t, nil, err)

	response, err = c.Get("http://localhost:8080/api/create")
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello", string(body))

}
