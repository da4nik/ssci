package ci

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/ssci/types"
)

const workspace = "workspace"

// Process processes notification
func Process(data types.Notificatable) {
	notification := data.Notification()
	workdir := filepath.Join(workspace, notification.Name)
	os.MkdirAll(workdir, os.ModePerm)

	log().Infof("Processing %s (%s)", notification.Name, notification.CloneURL)

	// TODO: #1 Catch output and save it to store with results

	if err := getSources(notification.CloneURL, workdir); err != nil {
		return
	}

	if err := runTests(workdir); err != nil {
		return
	}

	// TODO: #4 Build image
	// TODO: #5 Push image to docker repository
	// TODO: #6 Saving results to local storage
	// TODO: #7 Send resulting notifications
	// TODO: #9 Cleanup workspace
}

func getSources(url, workdir string) error {
	if _, err := os.Stat(filepath.Join(workdir, ".git")); os.IsNotExist(err) {
		return cloneRepo(url, workdir)
	}
	return updateCode(workdir)
}

func cloneRepo(url, workdir string) error {
	args := []string{"clone", url, workdir}
	cmd := exec.Command("git", args...)

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
	cmd := exec.Command("git", args...)
	cmd.Dir = workdir

	out, err := cmd.Output()
	if err != nil {
		log().Errorf("Unable to update repo: %v", err)
		return err
	}

	log().Debugf("%s updated. %s", workdir, out)

	return nil
}

func runTests(workdir string) error {
	args := []string{"test"}
	cmd := exec.Command("make", args...)
	cmd.Dir = workdir

	out, err := cmd.Output()
	if err != nil {
		log().Errorf("Unable run tests: %v", err)
		return err
	}

	log().Debugf("%s tests passed. %s", workdir, out)

	return nil
}

func log() *logrus.Entry {
	return logrus.WithField("module", "ci")
}
