package url

import (
	"strings"
	"fmt"
)

func ChangePathToUrl(path string) string {
	fmt.Println(path, strings.Replace(path, "\\", "/", -1))
	return strings.Replace(path, "\\", "/", -1)
}
