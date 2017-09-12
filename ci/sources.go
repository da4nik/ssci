package ci

import (
	"os"
	"os/exec"
	"path/filepath"
)

// TODO: #14 automatically locate git binary
const gitCommand = "/usr/bin/git"

func getSources(url, workdir string) error {
	if _, err := os.Stat(filepath.Join(workdir, ".git")); os.IsNotExist(err) {
		return cloneRepo(url, workdir)
	}
	return updateCode(workdir)
}

func cloneRepo(url, workdir string) error {
	args := []string{"clone", url, workdir}
	cmd := exec.Command(gitCommand, args...)

	out, err := cmd.Output()
	if err != nil {
		log().Errorf("Unable to clone repo: %v", err)
		return err
	}

	log().Debugf("%s cloned. %s", url, out)

	return nil
}

func updateCode(workdir string) error {
	args := []string{"pull"}
	cmd := exec.Command(gitCommand, args...)
	cmd.Dir = workdir

	out, err := cmd.Output()
	if err != nil {
		log().Errorf("Unable to update repo: %v", err)
		return err
	}

	log().Debugf("%s updated. %s", workdir, out)

	return nil
}
