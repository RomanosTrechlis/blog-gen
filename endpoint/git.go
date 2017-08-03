package endpoint

import (
	"fmt"
	"os/exec"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
	"strings"
)

// gitEndpoint is the git endpoint object
type gitEndpoint struct{}

// newGitEndpoint creates a new GitEndpoint
func newGitEndpoint() (e Endpoint) {
	return &gitEndpoint{}
}

// Upload uploads the site to a git repository
// todo: push fails
func (ds *gitEndpoint) Upload(to string) (err error) {
	fmt.Println("Uploading Site...")
	path := config.SiteInfo.DestFolder
	dest := config.SiteInfo.DestFolder + "_upload"
	err = fs.CreateFolderIfNotExist(dest)
	if err != nil {
		return err
	}
	err = fs.ClearFolder(dest)
	if err != nil {
		return err
	}

	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", path, err)
	}

	url, err := createUrlWithCred(to)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	remoteArgs := []string{"remote", "add", "origin", url}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = dest
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("error creating upload folder %s: %v", dest, err)
	}
	err = fs.CopyDir(config.SiteInfo.DestFolder, dest)
	if err != nil {
		return fmt.Errorf("error copying generated folder %s to upload folder %s: %v",
			config.SiteInfo.DestFolder, dest, err)
	}

	addArgs := []string{"add", "."}
	cmd = exec.Command(cmdName, addArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding files to commit: %v", err)
	}

	commitArgs := []string{"commit", "-m", "auto commit"}
	cmd = exec.Command(cmdName, commitArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error commiting files: %v", err)
	}

	pushArgs := []string{"push", "origin", "master"}
	cmd = exec.Command(cmdName, pushArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing to remote %s: %v", to, err)
	}
	fmt.Println("Upload Complete.")
	return nil
}

func createUrlWithCred(to string) (url string, err error) {
	t := strings.Split(to, "://")
	if len(t) != 2 {
		return "", fmt.Errorf("couldn't process git url")
	}
	p := strings.Replace(config.SiteInfo.Upload.Password, "@", "%40", 5)
	return t[0] + "://" + config.SiteInfo.Upload.Username + ":" + p + "@" + t[1], nil
}
