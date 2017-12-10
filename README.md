# Go Splitter

Library for convenient HTTP route splitting

## What for ?

Go provide out of the box route handling, but it is not possible to 
'split' route part and handle entry part with custom type - now it's possible

## Hot to use

```golang
  err := gosplitter.Match("/your/route", http.Handle|gosplitter.HandlerFunc|interface{})
```

You can pass handler or another structure to chain multiple routes

No additional params, easy as it does!

# TODO:

- More tests
- Benchmarking
- Clear code
- Comment-out code