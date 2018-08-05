package url_test

import (
	"testing"

	"github.com/RomanosTrechlis/blog-gen/util/url"
)

func TestChangePathToUrl(t *testing.T) {
	tests := []struct {
		path, url string
	}{
		{"folder\\folder", "folder/folder"},
		{"folder/folder", "folder/folder"},
		{"folder", "folder"},
	}

	for _, tt := range tests {
		u := url.ChangePathToUrl(tt.path)
		if u != tt.url {
			t.Errorf("expected '%s', got '%s'", tt.url, u)
		}
	}
}
