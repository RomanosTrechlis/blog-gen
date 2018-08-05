package fs_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RomanosTrechlis/blog-gen/util/fs"
)

func TestCreateFolderIfNotExist(t *testing.T) {
	var tests = []struct {
		folder       string
		desc         string
		err          bool
		removeFolder bool
	}{
		{filepath.Join("testdata", "existing"), "folder exists", false, false},
		{filepath.Join("testdata", "notExisting"), "folder exists", false, true},
	}

	for _, tt := range tests {
		err := fs.CreateFolderIfNotExist(tt.folder)
		if err != nil && !tt.err {
			t.Errorf("expected no error, got %v", err)
		}
		if err == nil && tt.err {
			t.Errorf("expected error, got no error")
		}

		if tt.removeFolder {
			os.Remove(tt.folder)
		}
	}
}

func TestCopyDir(t *testing.T) {
	var tests = []struct {
		from         string
		to           string
		desc         string
		err          bool
		removeFolder bool
	}{
		{filepath.Join("testdata", "existing"), filepath.Join("testdata", "copy"), "", false, true},
		{filepath.Join("testdata", "nofolder"), filepath.Join("testdata", "copy"), "", true, true},
	}

	for _, tt := range tests {
		err := fs.CopyDir(tt.from, tt.to)
		if err != nil && !tt.err {
			t.Errorf("expected no error, got %v", err)
		}
		if err == nil && tt.err {
			t.Errorf("expected error, got no error")
		}
		if err != nil {
			continue
		}

		if tt.removeFolder {
			os.RemoveAll(tt.to)
		}
	}
}

func TestGetFilenameFrom(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"folder\\file.txt", "file.txt"},
		{"folder/file.txt", "folder/file.txt"}, // shouldn't be like that but it's ok for now
		{"folder\\folder\\file.txt", "file.txt"},
		{"file.txt", "file.txt"},
		{"\\file.txt", "file.txt"},
	}
	if runtime.GOOS != "windows" {
		tests = []struct {
			path   string
			result string
		}{
			{"folder/file.txt", "file.txt"},
			{"folder\\file.txt", "folder/file.txt"}, // shouldn't be like that but it's ok for now
			{"folder/folder/file.txt", "file.txt"},
			{"file.txt", "file.txt"},
			{"/file.txt", "file.txt"},
		}
	}

	for _, tt := range tests {
		r := fs.GetFilenameFrom(tt.path)
		if r != tt.result {
			t.Errorf("expected '%s', got '%s' (%s)", tt.result, r, tt.path)
		}
	}
}

func TestGetFolderNameFrom(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"folder\\file.txt", "folder"},
		{"folder/file.txt", ""},
		{"folder\\folder\\file.txt", "folder\\folder"},
		{"file.txt", ""},
		{"\\file.txt", ""},
	}
	if runtime.GOOS != "windows" {
		tests = []struct {
			path   string
			result string
		}{
			{"folder/file.txt", "folder"},
			{"folder\\file.txt", ""},
			{"folder/folder/file.txt", "folder\\folder"},
			{"file.txt", ""},
			{"/file.txt", ""},
		}
	}

	for _, tt := range tests {
		r := fs.GetFolderNameFrom(tt.path)
		if r != tt.result {
			t.Errorf("expected '%s', got '%s' (%s)", tt.result, r, tt.path)
		}
	}
}
