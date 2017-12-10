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
	var result = unstar(unbrace(unslash(f.Name())))

	// fmt.Printf("RETURNING %s\n", result)

	return result
}

// GetFunctionName func
func GetFunctionName(myvar interface{}) string {
	var result string
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		result = t.Elem().Name()
	} else {
		result = t.Name()
	}

	return result
}

func unslash(s string) string {
	r := strings.Split(s, "/")
	return r[len(r)-1]
}

func undot(s string) string {
	r := strings.Split(s, ".")
	return r[len(r)-1]
}

func unbrace(s string) string {
	var result string
	result = strings.Replace(s, "(", "", -1)
	result = strings.Replace(result, ")", "", -1)
	return result
}

func unstar(s string) string {
	return strings.Replace(s, "*", "", -1)
}
