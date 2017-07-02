package datasource

import (
	"fmt"

	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// LocalDataSource is the local data source object
type LocalDataSource struct{}

// NewLocalDataSource creates a new LocalDataSource
func NewLocalDataSource() DataSource {
	return &LocalDataSource{}
}

// Fetch creates the output folder, clears it and copies the local folder there
func (ds *LocalDataSource) Fetch(from, to string) ([]string, error) {
	fmt.Printf("Fetching data from %s into %s...\n", from, to)
	err := fs.CreateFolderIfNotExist(to)
	if err != nil {
		return nil, err
	}
	err = fs.ClearFolder(to)
	if err != nil {
		return nil, err
	}
	err = fs.CopyDir(from, to)
	if err != nil {
		return nil, err
	}
	dirs, err := fs.GetContentFolders(to)
	if err != nil {
		return nil, err
	}
	fmt.Print("Fetching complete.\n")
	return dirs, nil
}
