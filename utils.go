package gosplitter

import (
	"reflect"
	"runtime"
	"strings"
)

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

	return unslash(f.Name())
}

// GetFunctionName func
func GetFunctionName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func unslash(s string) string {
	r := strings.Split(s, "/")
	return r[len(r)-1]
}

func undot(s string) string {
	r := strings.Split(s, ".")
	return r[len(r)-1]
}
