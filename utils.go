package swift

import (
	"fmt"
	"regexp"
	"strings"
)

func Prepend[T any](slice []T, elems ...T) []T {
	return append(elems, slice...)
}

func BuildAndValidatePath(path string) string {
	p := normalizePath(path)

	if !validPath(p) {
		panic(fmt.Sprintf("Invalid Path: '%s' - Original Path: '%s'", p, path))
	}

	return p
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if strings.HasSuffix(path, "/") && path != "/" {
		path = path[:len(path)-1]
	}

	return path
}

func validPath(path string) bool {
	pattern := `^\/[^\/].*[^\/]$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(path)
}
