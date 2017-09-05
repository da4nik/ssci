package ci

import (
	"github.com/Sirupsen/logrus"
	"github.com/da4nik/ssci/types"
)

// Process processes notification
func Process(data types.Notificatable) {
	notification := data.Notification()

	log().Infof("Processing %s (%s)", notification.Name, notification.CloneURL)

	// TODO: #1 Catch output and save it to store with results

	// TODO: #2 Clone repository
	// TODO: #3 Run tests
	// TODO: #4 Build image
	// TODO: #5 Push image to docker repository
	// TODO: #6 Saving results to local storage
	// TODO: #7 Send resulting notifications
}

func log() *logrus.Entry {
	return logrus.WithField("module", "ci")
}
