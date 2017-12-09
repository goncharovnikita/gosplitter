package gosplitter_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goncharovnikita/gosplitter"
)

type UsersRouterPoint struct{}
type CreateUsersHandler struct{}

func (c CreateUsersHandler) Handle(w http.ResponseWriter, r *http.Request) {}
func (u UsersRouterPoint) ServeUsers() {
	var h CreateUsersHandler
	gosplitter.Match("/create", h)
}
func TestIntegration(t *testing.T) {
	var u UsersRouterPoint
	var err error
	var rootCaller = "github.com/goncharovnikita/gosplitter_test.TestIntegration"
	// var serveUsersCaller = "github.com/goncharovnikita/gosplitter_test.ServeUsers"

	err = gosplitter.Match("/api", u)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, gosplitter.RegisteredPatterns[rootCaller].Children != nil)
	assert.Equal(t, true, gosplitter.RegisteredPatterns[rootCaller].Children["/api"] != nil)
	// assert.Equal(t, true, gosplitter.RegisteredPatterns[serveUsersCaller].Children != nil)
	// assert.Equal(t, true, gosplitter.RegisteredPatterns[serveUsersCaller].Children["/create"] != nil)
	t.Logf("%v", gosplitter.RegisteredPatterns)
}
