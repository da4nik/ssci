package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/ssci/store"
	"github.com/da4nik/ssci/types"
)

const workspace = "workspace"

// Process processes notification
func Process(data types.Notificatable) {
	buildStart := time.Now()
	notification := data.Notification()
	workdir := filepath.Join(workspace, notification.Name)
	os.MkdirAll(workdir, os.ModePerm)

	var project *types.Project
	project, err := store.LoadProject(notification.Name)
	if err != nil {
		project = store.NewProject(notification.Name, notification.CloneURL)
	}

	build := store.NewBuild(project)
	imageTag := fmt.Sprintf("%s:%d", notification.Name, build.ID)

	logp := logrus.WithFields(logrus.Fields{
		"project": project.Name,
		"build#":  build.ID,
		"module":  "ci",
	})

	logp.Infof("Processing #%d %s (%s)", build.ID, notification.Name, notification.CloneURL)

	// TODO: #1 Catch output and save it to store with results

	start := time.Now()
	logp.Debugf("Getting sources for %s", notification.CloneURL)
	if err := getSources(notification.CloneURL, workdir); err != nil {
		logp.Errorf("Get sources error: %v", err)
		return
	}
	logp.Debugf("Got sources for %s", notification.CloneURL, time.Since(start))

	start = time.Now()
	logp.Debugf("Running tests %s (%s)", notification.Name, time.Since(start))
	if err := runTests(workdir); err != nil {
		logp.Errorf("Run tests error: %v", err)
		return
	}
	logp.Debugf("Tests are passed (%s)", time.Since(start))

	start = time.Now()
	logp.Debugf("Building image \"%s\"", imageTag)
	if err := buildImage(imageTag, workdir); err != nil {
		logp.Errorf("Build image error: %v", err)
		return
	}
	logp.Debugf("Image \"%s\" built (%s)", notification.Name, time.Since(start))

	// TODO: #5 Push image to docker repository
	// TODO: #6 Saving results to local storage
	// TODO: #7 Send resulting notifications

	start = time.Now()
	logp.Debugf("Removing workspace: %s", workdir)
	if err := cleanupWorkspace(workdir); err != nil {
		logp.Errorf("Error cleaning up workspace: %v", err)
		return
	}
	logp.Debugf("%s cleaned up (%s)", workdir, time.Since(start))

	logp.Infof("Build of %s is done (%s)", notification.Name, time.Since(buildStart))

	build.Duration = string(time.Since(buildStart))
	if err := store.SaveProject(project); err != nil {
		logp.Errorf("Error saving project: %v", err)
	}

	start = time.Now()
	logp.Debugf("Starting image %s service update for %s", imageTag, notification.Name)
	if err := updateSwarmServiceImage(notification.Name, imageTag); err != nil {
		logp.Errorf("Error updating service image: %v", err)
		return
	}
	logp.Debugf("Service %s image updated (%s)", notification.Name, time.Since(start))

}

func log() *logrus.Entry {
	return logrus.WithField("module", "ci")
}
