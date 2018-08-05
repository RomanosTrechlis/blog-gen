package url

import (
	"strings"
)

func ChangePathToUrl(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
