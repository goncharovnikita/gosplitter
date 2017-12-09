package gosplitter

import (
	"fmt"
)

// NotRegisteredPatternError error
type NotRegisteredPatternError struct {
	pattern string
}

func (n *NotRegisteredPatternError) Error() string {
	return fmt.Sprintf("pattern %v not registered", n.pattern)
}

// InvalidHandlerError error
type InvalidHandlerError struct {
	handler interface{}
}

func (e InvalidHandlerError) Error() string {
	return fmt.Sprintf(
		"invalid handler: %v\nAvailable handler types: http.Hander, gosplitter.HandlerFunc, gosplitter.RouterPoint",
		e.handler,
	)
}

// PatternAlreadyRegisteredError error
type PatternAlreadyRegisteredError struct {
	name string
}

func (e PatternAlreadyRegisteredError) Error() string {
	return fmt.Sprintf("pattern %s already registered", e.name)
}
