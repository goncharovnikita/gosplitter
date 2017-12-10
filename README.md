# Go Splitter

[![Build Status](https://travis-ci.org/goncharovnikita/gosplitter.svg?branch=master)](https://travis-ci.org/goncharovnikita/gosplitter) 
[![GoDoc](https://godoc.org/github.com/goncharovnikita/gosplitter?status.svg)](https://godoc.org/github.com/goncharovnikita/gosplitter)
[![Go Report Card](https://goreportcard.com/badge/github.com/goncharovnikita/gosplitter)](https://goreportcard.com/report/github.com/goncharovnikita/gosplitter)

Library for convenient HTTP route splitting

## What for ?

Go provide out of the box route handling, but it is not possible to 
'split' route part and handle entry part with custom type - now it's possible

## Hot to use

```golang
  gosplitter.Match("/your/path", *http.ServeMux, http.Handler|http.HandleFunc|interface{})
```

You can pass handler or another structure to chain multiple routes

----

## Examples

```
/**
* Specify handle types
 */
type APIV1Handler struct {
	mux *http.ServeMux
}

type ColorsHandler struct {
	mux *http.ServeMux
}

/**
* Start - binds parent to children
 */
func (a *APIV1Handler) Start() {
	var colorsHandler = ColorsHandler{
		mux: a.mux,
	}
	gosplitter.Match("/ping", a.mux, a.HandlePing())
	gosplitter.Match("/colors", a.mux, colorsHandler)
	colorsHandler.Start()
}

/**
* Simple http handler function
 */
func (a *APIV1Handler) HandlePing() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func (c *ColorsHandler) Start() {
	gosplitter.Match("/black", c.mux, c.HandleBlack())
}

func (c *ColorsHandler) HandleBlack() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("#000000"))
	}
}

func main() {
  var mux = http.NewServeMux()
	var apiV1 = APIV1Handler{
		mux: mux,
	}

  /**
  * bind api handler to root
   */
  gosplitter.Match("/api/v1", mux, apiV1)
  /**
  * start api handler
   */

  apiV1.Start()
 
}
```

----

# TODO:

- More tests
- Benchmarking
- Clear code
- Comment-out code