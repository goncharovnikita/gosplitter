package gosplitter

import (
	"net/http"
	"runtime"
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
func Match(url string, h interface{}) error {
	var handler http.Handler
	var handlerFunc HandlerFunc
	var mux = http.NewServeMux()

	switch v := h.(type) {
	case http.Handler:
		handler = v
		RegisterHandler(url, CallerContext(), handler, mux)
		break
	case HandlerFunc:
		handlerFunc = v
		RegisterHandler(url, CallerContext(), handlerFunc, mux)
		break
	default:
		RegisterRouterPoint(url, CallerContext())
	}

	return nil
}

// RegisterHandler func
func RegisterHandler(url string, caller string, h interface{}, mux *http.ServeMux) error {
	var handler http.Handler
	var handlerFunc HandlerFunc
	switch v := h.(type) {
	case http.Handler:
		handler = v
		break
	case HandlerFunc:
		handlerFunc = v
		break
	}
	point := RegisteredPatterns[caller]
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

	var absURL, err = GetAbsoluteURL(url, caller)
	if err != nil {
		return err
	}
	if point.Children == nil {
		rMap := make(map[string]*RouterPoint)
		point.Children = rMap
		if handler != nil {
			point.Children[url] = &RouterPoint{
				URL:     absURL,
				Handler: handler,
			}
			mux.Handle(absURL, handler)
			return nil
		}
		point.Children[url] = &RouterPoint{
			URL:         absURL,
			HandlerFunc: handlerFunc,
		}
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
		mux.Handle(absURL, handler)
		return nil
	}
	point.Children[url] = &RouterPoint{
		URL:         absURL,
		HandlerFunc: handlerFunc,
	}
	mux.HandleFunc(absURL, func(w http.ResponseWriter, r *http.Request) {
		handlerFunc.Handle(w, r)
	})
	return nil
}

// RegisterRouterPoint func
func RegisterRouterPoint(url string, caller string) error {
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

// GetAbsoluteURL func
func GetAbsoluteURL(url string, caller string) (string, error) {
	point := RegisteredPatterns[caller]
	if point == nil {
		return "", &NotRegisteredPatternError{
			pattern: url,
		}
	}
	return point.URL + url, nil
}

// CallerContext func
func CallerContext() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	f := runtime.FuncForPC(fpcs[0] - 1)
	if f == nil {
		return "n/a"
	}

	return f.Name()
}
