# Go Splitter

[![Build Status](https://travis-ci.org/goncharovnikita/gosplitter.svg?branch=master)](https://travis-ci.org/goncharovnikita/gosplitter) 
[![GoDoc](https://godoc.org/github.com/goncharovnikita/gosplitter?status.svg)](https://godoc.org/github.com/goncharovnikita/gosplitter)

Library for convenient HTTP route splitting

## What for ?

Go provide out of the box route handling, but it is not possible to 
'split' route part and handle entry part with custom type - now it's possible

## Hot to use

```golang
  err := gosplitter.Match("/your/route", *http.serveMux, http.Handle|gosplitter.HandlerFunc|interface{})
```

You can pass handler or another structure to chain multiple routes

No additional params, easy as it does!

# TODO:

- More tests
- Benchmarking
- Clear code
- Comment-out code