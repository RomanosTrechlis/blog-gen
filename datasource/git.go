package datasource

import (
	"fmt"
	"os/exec"

	"github.com/RomanosTrechlis/blog-gen/util/fs"
)

// gitDataSource is the git data source object
type gitDataSource struct{}

// newGitDataSource creates a new GitDataSource
func newGitDataSource() (ds DataSource) {
	return &gitDataSource{}
}

// Fetch creates the output folder, clears it and clones the repository there
func (ds *gitDataSource) Fetch(from, to string) (dirs []string, err error) {
	fmt.Printf("Fetching data from %s into %s...\n", from, to)
	err = fs.CreateFolderIfNotExist(to)
	if err != nil {
		return nil, err
	}
	err = fs.ClearFolder(to)
	if err != nil {
		return nil, err
	}
	err = cloneRepo(to, from)
	if err != nil {
		return nil, err
	}
	dirs, err = fs.GetContentFolders(to)
	if err != nil {
		return nil, err
	}
	fmt.Println("Fetching complete.")
	return dirs, nil
}

func cloneRepo(path, repositoryURL string) (err error) {
	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", path, err)
	}

	remoteArgs := []string{"remote", "add", "origin", repositoryURL}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting remote %s: %v", repositoryURL, err)
	}

	pullArgs := []string{"pull", "origin", "master"}
	cmd = exec.Command(cmdName, pullArgs...)
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error pulling master at %s: %v", path, err)
	}
	return nil
}
