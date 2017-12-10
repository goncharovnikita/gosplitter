package gosplitter

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc interface
type HandlerFunc interface {
	Handle(http.ResponseWriter, *http.Request)
}

// RouterPoint type
type RouterPoint struct {
	URL         string
	Handler     http.Handler
	HandlerFunc HandlerFunc
	Children    map[string]*RouterPoint
}

// RegisteredPatterns type
var RegisteredPatterns = make(map[string]*RouterPoint)

// Match matches route
func Match(url string, mux *http.ServeMux, h interface{}) error {
	fmt.Printf("%s\n", url)
	var handler http.Handler
	var handlerFunc HandlerFunc
	callerContext := CallerContext()
	fmt.Printf("Caller context: %s\n", callerContext)

	switch v := h.(type) {
	case http.Handler:
		handler = v
		RegisterHandler(url, callerContext, handler, mux)
		break
	case HandlerFunc:
		handlerFunc = v
		RegisterHandler(url, callerContext, handlerFunc, mux)
		break
	default:
		RegisterRouterPoint(url, h)
	}

	return nil
}

// RegisterHandler func
func RegisterHandler(url string, caller string, h interface{}, mux *http.ServeMux) error {
	var handler http.Handler
	var handlerFunc HandlerFunc
	var callerPoints = strings.Split(caller, ".")
	var point *RouterPoint
	switch v := h.(type) {
	case http.Handler:
		handler = v
		break
	case HandlerFunc:
		handlerFunc = v
		break
	}

	for i := len(callerPoints) - 1; i >= 0; i-- {
		if v := RegisteredPatterns[callerPoints[i]]; v != nil {
			point = v
			break
		}
	}

	if point == nil {
		if handler != nil {
			RegisteredPatterns[caller] = &RouterPoint{
				URL:     url,
				Handler: handler,
			}
			mux.Handle(url, handler)
			return nil
		}
		RegisteredPatterns[caller] = &RouterPoint{
			URL:         url,
			HandlerFunc: handlerFunc,
		}
		mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			handlerFunc.Handle(w, r)
		})
		return nil
	}

	if point.Handler != nil {
		return PatternAlreadyRegisteredError{
			name: url,
		}
	}

	var absURL = point.URL + url
	if point.Children == nil {
		rMap := make(map[string]*RouterPoint)
		point.Children = rMap
		if handler != nil {
			point.Children[url] = &RouterPoint{
				URL:     absURL,
				Handler: handler,
			}
			fmt.Printf("Handle %s\n", absURL)
			mux.Handle(absURL, handler)
			return nil
		}
		point.Children[url] = &RouterPoint{
			URL:         absURL,
			HandlerFunc: handlerFunc,
		}
		fmt.Printf("Handle %s\n", absURL)
		mux.HandleFunc(absURL, func(w http.ResponseWriter, r *http.Request) {
			handlerFunc.Handle(w, r)
		})
		return nil
	}

	if point.Children[absURL] != nil {
		return PatternAlreadyRegisteredError{
			name: absURL,
		}
	}

	if handler != nil {
		point.Children[url] = &RouterPoint{
			URL:     absURL,
			Handler: handler,
		}
		fmt.Printf("Handle %s\n", absURL)
		mux.Handle(absURL, handler)
		return nil
	}
	point.Children[url] = &RouterPoint{
		URL:         absURL,
		HandlerFunc: handlerFunc,
	}
	fmt.Printf("Handle %s\n", absURL)
	mux.HandleFunc(absURL, func(w http.ResponseWriter, r *http.Request) {
		handlerFunc.Handle(w, r)
	})
	return nil
}

// RegisterRouterPoint func
func RegisterRouterPoint(url string, f interface{}) error {
	var caller = undot(GetFunctionName(f))
	fmt.Printf("register router point for %s\n", caller)
	var point = RegisteredPatterns[caller]
	if point == nil {
		// Point does not registered yet
		RegisteredPatterns[caller] = &RouterPoint{
			URL: url,
			Children: map[string]*RouterPoint{
				url: &RouterPoint{
					URL: url,
				},
			},
		}
		return nil
	}

	if point.Handler != nil {
		return PatternAlreadyRegisteredError{
			name: url,
		}
	}

	var absURL, err = GetAbsoluteURL(url, caller)
	if err != nil {
		return err
	}

	if point.Children == nil {
		point.Children = map[string]*RouterPoint{
			absURL: &RouterPoint{
				URL: url,
			},
		}
		return nil
	}

	if point.Children[absURL] != nil {
		return PatternAlreadyRegisteredError{
			name: absURL,
		}
	}

	point.Children[absURL] = &RouterPoint{
		URL: url,
	}
	return nil
}
