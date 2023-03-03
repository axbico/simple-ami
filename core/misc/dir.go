package misc

import (
	"net/http"
	"os"
	"regexp"
	"strings"
)

/// ///

const (
	PUBLIC_RELATIVE_PATH_DIRECTORY = "public/"
	PUBLIC_ALLOWED_EXTENSIONS      = "(ico|jpg|jpeg|png|html|txt)$"

	APPLICATION_RELATIVE_PATH_DIRECTORY = "application/"
	APPLICATION_ALLOWED_EXTENSIONS      = "(html|htm|css|js)$"
)

/// ///

func allowedFileAccess(request *http.Request, relativePathDirectory string, allowedExtensions string) (string, bool) {
	path := request.URL.Path

	if regexp.MustCompile(allowedExtensions).MatchString(path[strings.LastIndexByte(path, '.')+1:]) {
		filepath := relativePathDirectory + path
		if _, err := os.Stat(filepath); err == nil {
			return filepath, true
		}
	}

	return "", false
}

func AllowedFilePublic(request *http.Request) (string, bool) {
	return allowedFileAccess(
		request,
		PUBLIC_RELATIVE_PATH_DIRECTORY,
		PUBLIC_ALLOWED_EXTENSIONS,
	)
}

func AllowedFileApplication(request *http.Request) (string, bool) {
	return allowedFileAccess(
		request,
		APPLICATION_RELATIVE_PATH_DIRECTORY,
		APPLICATION_ALLOWED_EXTENSIONS,
	)
}

/// ///
