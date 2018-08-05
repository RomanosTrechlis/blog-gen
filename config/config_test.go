package config_test

import (
	"testing"
	"github.com/RomanosTrechlis/blog-gen/config"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		file string
		err bool
	} {
		{"testdata\\config.json", false},
		{"testdata\\nofile.json", true},
		{"testdata\\configFillValues.json", false},
	}

	for _, tt := range tests {
		s, err := config.New(tt.file)

		if err != nil && !tt.err {
			t.Fatalf("expected no error, got %v", err)
		}
		if err == nil && tt.err {
			t.Errorf("expected error, got no error")
		}
		if err != nil {
			continue
		}

		if s.TempFolder != "./tmp" {
			t.Errorf("expected templ folder to be '%s', got '%s'", "./tmp", s.TempFolder)
		}
		if s.DestFolder != "./public" {
			t.Errorf("expected destination folder to be '%s', got '%s'", "./public", s.DestFolder)
		}
		if s.ThemeFolder != "./static/" {
			t.Errorf("expected theme folder to be '%s', got '%s'", "./static", s.ThemeFolder)
		}
		if s.NumPostsFrontPage != 10 {
			t.Errorf("expected number of posts to be '10', got '%d'", s.NumPostsFrontPage)
		}
	}
}
