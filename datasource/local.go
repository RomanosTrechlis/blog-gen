package datasource

import (
	"fmt"

	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// localDataSource is the local data source object
type localDataSource struct{}

// newLocalDataSource creates a new LocalDataSource
func newLocalDataSource() (ds DataSource) {
	return &localDataSource{}
}

// Fetch creates the output folder, clears it and copies the local folder there
func (ds *localDataSource) Fetch(from, to string) (dirs []string, err error) {
	fmt.Printf("Fetching data from %s into %s...\n", from, to)
	err = fs.CreateFolderIfNotExist(to)
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
	dirs, err = fs.GetContentFolders(to)
	if err != nil {
		return nil, err
	}
	fmt.Print("Fetching complete.\n")
	return dirs, nil
}
