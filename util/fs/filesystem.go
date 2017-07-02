package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", src, err)
	}
	defer in.Close()
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
	return nil
}

func GetContentFolders(path string) ([]string, error) {
	var result []string
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
			result = append(result, fmt.Sprintf("%s/%s", path, file.Name()))
		}
	}
	return result, nil
}

func CopyDir(source, path string) (err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return nil
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		src := fmt.Sprintf("%s/%s", source, file.Name())
		dst := fmt.Sprintf("%s/%s", path, file.Name())
		if file.IsDir() {
			err := os.Mkdir(dst, os.ModePerm)
			if err != nil {
				return err
			}
			err = CopyDir(src, dst)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateFolderIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("error accessing directory %s: %v", path, err)
		}
	}
	return nil
}

func ClearFolder(path string) error {
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
