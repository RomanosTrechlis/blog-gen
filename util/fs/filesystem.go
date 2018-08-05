package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

func CopyFile(src, dest string) (err error) {
	file := GetFilenameFrom(src)
	dst := filepath.Join(dest, file)
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", src, err)
	}
	defer in.Close()

	if dest != "" {
		err = CreateFolderIfNotExist(dest)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", dst, err)
	}

	defer func() {
		e := out.Close()
		if e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", dst, err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", dst, err)
	}
	return err
}

func GetContentFolders(path string) (result []string, err error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading contents of directory %s: %v", path, err)
	}
	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}
	return result, nil
}

func CopyDir(source, dest string) (err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		src := filepath.Join(source, file.Name())
		if !file.IsDir() {
			err = CopyFile(src, dest)
			if err != nil {
				return err
			}
			return nil
		}

		dst := filepath.Join(dest, file.Name())
		err := CreateFolderIfNotExist(dst)
		if err != nil {
			return err
		}
		err = CopyDir(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateFolderIfNotExist(path string) (err error) {
	_, err = os.Stat(path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error accessing directory %s: %v", path, err)
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %v", path, err)
	}
	return nil
}

func ClearFolder(path string) (err error) {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return fmt.Errorf("error clearing file %s: %v", name, err)
		}
	}
	return nil
}

func GetFolderNameFrom(path string) string {
	separator := GetSeparator()

	i := strings.LastIndex(path, separator) + 1
	if i >= len(path) {
		return path
	}

	if i -  1 >= 0 {
		return path[: i - 1]
	}

	return path[:i]
}

func GetFilenameFrom(path string) string {
	separator := GetSeparator()
	i := strings.LastIndex(path, separator) + 1
	if i >= len(path) {
		return ""
	}
	f := path[i:]
	if !strings.Contains(f, ".") {
		return ""
	}
	return f
}

func GetSeparator() string {
	separator := "/"
	osWin := runtime.GOOS == "windows"
	if osWin {
		separator = "\\"
	}
	return separator
}
