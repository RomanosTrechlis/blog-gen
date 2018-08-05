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
		{filepath.Join("testdata", "notExisting"), "folder doesn't exist", false, true},
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
		{"", ""},
		{"folder\\file", ""},
	}
	if runtime.GOOS != "windows" {
		tests = []struct {
			path   string
			result string
		}{
			{"folder/file.txt", "file.txt"},
			{"folder\\file.txt", "folder\\file.txt"}, // shouldn't be like that but it's ok for now
			{"folder/folder/file.txt", "file.txt"},
			{"file.txt", "file.txt"},
			{"/file.txt", "file.txt"},
			{"", ""},
			{"folder/file", ""},
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
		{"", ""},
	}
	if runtime.GOOS != "windows" {
		tests = []struct {
			path   string
			result string
		}{
			{"folder/file.txt", "folder"},
			{"folder\\file.txt", ""},
			{"folder/folder/file.txt", "folder/folder"},
			{"file.txt", ""},
			{"/file.txt", ""},
			{"", ""},
		}
	}

	for _, tt := range tests {
		r := fs.GetFolderNameFrom(tt.path)
		if r != tt.result {
			t.Errorf("expected '%s', got '%s' (%s)", tt.result, r, tt.path)
		}
	}
}

func TestGetContentFolders(t *testing.T) {
	var tests = []struct {
		folder string
		desc   string
		err    bool
		num    int
	}{
		{filepath.Join("testdata", "existing"), "folder exists and has 1 file", false, 0},
		{filepath.Join("testdata", "notExisting"), "folder doesn't exist", true, 0},
		{filepath.Join("testdata"), "folder exists and has 3 folders", false, 3},
	}

	for _, tt := range tests {
		c, err := fs.GetContentFolders(tt.folder)
		if err != nil && !tt.err {
			t.Errorf("expected no error, got '%v'", err)
		}
		if err == nil && tt.err {
			t.Error("expected error, got no error")
		}
		if err != nil {
			continue
		}

		if len(c) != tt.num {
			t.Errorf("expected %d number of files, got %d", tt.num, len(c))
		}
	}
}

func TestClearFolder(t *testing.T) {
	var tests = []struct {
		folder       string
		err          bool
		createFolder bool
		removeFolder bool
	}{
		{filepath.Join("testdata", "tempFolder"), false, true,true},
		{filepath.Join("testdata", "tempFolder"), true, false,false},
	}

	cp := filepath.Join("testdata", "existing")
	for _, tt := range tests {

		if tt.createFolder {
			err := fs.CopyDir(cp, tt.folder)
			if err != nil {
				t.Errorf("failed to initialize test: copy of dir '%s' to dir '%s' failed: stopping test", cp, tt.folder)
				continue
			}
		}

		err := fs.ClearFolder(tt.folder)
		if err != nil && !tt.err {
			t.Errorf("expected no error, got '%v'", err)
		}
		if err == nil && tt.err {
			t.Error("expected error, got no error")
		}
		if err != nil {
			continue
		}

		dir, err := os.Open(tt.folder)
		if err != nil {
			t.Errorf("error accessing directory %s: %v", tt.folder, err)
			continue
		}
		defer dir.Close()
		files, err := dir.Readdir(-1)
		if err != nil {
			t.Errorf("error reading contents of directory %s: %v", tt.folder, err)
			continue
		}

		if len(files) != 0 {
			t.Errorf("expecting 0 files in directory %s, got %d", tt.folder, len(files))
		}

		if tt.removeFolder {
			os.RemoveAll(tt.folder)
		}
	}
}
