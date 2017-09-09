package ci

import (
	"os"
	"path/filepath"
	"time"

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

	start := time.Now()
	log().Debugf("Getting sources for %s", notification.CloneURL)
	if err := getSources(notification.CloneURL, workdir); err != nil {
		log().Errorf("Get sources error: %v", err)
		return
	}
	log().Debugf("Got sources for %s", notification.CloneURL, time.Since(start))

	start = time.Now()
	log().Debugf("Running tests %s (%s)", notification.Name, time.Since(start))
	if err := runTests(workdir); err != nil {
		log().Errorf("Run tests error: %v", err)
		return
	}
	log().Debugf("Tests are passed (%s)", time.Since(start))

	// TODO: #10 Add version number from event or from build number for example
	start = time.Now()
	log().Debugf("Building image: %s", notification.Name)
	if err := buildImage(notification.Name, workdir); err != nil {
		log().Errorf("Build image error: %v", err)
		return
	}
	log().Debugf("Image \"%s\" built (%s)", notification.Name, time.Since(start))

	// TODO: #5 Push image to docker repository
	// TODO: #6 Saving results to local storage
	// TODO: #7 Send resulting notifications
	// TODO: #9 Cleanup workspace
}

func log() *logrus.Entry {
	return logrus.WithField("module", "ci")
}
