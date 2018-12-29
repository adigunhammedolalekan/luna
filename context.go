package luna

import (
	"fmt"
	"strings"
)

type Context struct {
	Path string
	Vars map[string]interface{}
	Data interface{}
}

// MatchRoute checks if path matches the patten supplied.
// E.g pattern=/rooms/{id}/messages will match path=/rooms/12/messages
func MatchRoute(pattern, path string) bool {

	namedParamsCount := 0
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for _, v := range patternParts {

		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			namedParamsCount++
		}
	}

	unMatchedCount := 0
	for i, v := range pathParts {

		if patternParts[i] != v {
			unMatchedCount++
		}
	}

	return unMatchedCount == namedParamsCount
}

// ExtractParams extracts data from path parameter.
func ExtractParams(template, path string) (map[string]interface{}, error) {

	data := make(map[string]interface{})

	templateParts := strings.Split(template, "/")
	pathParts := strings.Split(path, "/")

	if len(templateParts) != len(pathParts) {

		return make(map[string]interface{}), fmt.Errorf("Template and path does not match! %s %s", template, path)
	}

	for i, p := range templateParts {

		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			name := p[1 : len(p)-1]
			data[name] = pathParts[i]
		}
	}

	return data, nil
}
