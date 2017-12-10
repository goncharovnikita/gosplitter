package gosplitter

import (
	"net/http"
	"strings"
)

// RouterPoint type
type RouterPoint struct {
	URL         string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
}

// RegisteredPatterns type
var RegisteredPatterns = make(map[string]*RouterPoint)

// Match matches route
func Match(url string, mux *http.ServeMux, h interface{}) {
	callerContext := CallerContext()
	switch h.(type) {
	case http.Handler:
		RegisterHandler(url, callerContext, h, mux)
		break
	case func(http.ResponseWriter, *http.Request):
		RegisterHandler(url, callerContext, h, mux)
		break
	default:
		RegisterRouterPoint(callerContext, url, h)
	}
}

// RegisterHandler func
func RegisterHandler(url string, caller string, h interface{}, mux *http.ServeMux) {
	var handler http.Handler
	var handlerFunc http.HandlerFunc
	var callerPoints = strings.Split(caller, ".")
	var point *RouterPoint
	switch v := h.(type) {
	case http.Handler:
		handler = v
		break
	case func(http.ResponseWriter, *http.Request):
		handlerFunc = v
		break
	}

	for i := len(callerPoints) - 1; i >= 0; i-- {
		if v := RegisteredPatterns[callerPoints[i]]; v != nil {
			point = v
			break
		}
	}

	if point != nil {
		url = point.URL + url
	}

	if handler == nil {
		mux.HandleFunc(url, handlerFunc)
	} else {
		mux.Handle(url, handler)
	}

	// fmt.Printf("HANDLE %s\nCALLER %s\n", url, caller)

	RegisteredPatterns[caller] = &RouterPoint{
		URL:         url,
		Handler:     handler,
		HandlerFunc: handlerFunc,
	}
}

// RegisterRouterPoint func
func RegisterRouterPoint(caller string, url string, f interface{}) {
	var funcName = undot(GetFunctionName(f))
	// fmt.Printf("REGISTER POINT CALLER %s\n", caller)
	var point *RouterPoint
	var callerPoints = strings.Split(caller, ".")

	for i := len(callerPoints) - 1; i >= 0; i-- {
		if v := RegisteredPatterns[callerPoints[i]]; v != nil {
			point = v
			break
		}
	}

	if point != nil {
		url = point.URL + url
	}

	RegisteredPatterns[funcName] = &RouterPoint{
		URL: url,
	}
}
