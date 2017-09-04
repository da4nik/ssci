package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/da4nik/ssci/webhooks"
)

var version string
var buildTime string

var logFile *os.File

func initLogger(logFileName string) {
	// Setting up logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)

	log.SetOutput(os.Stdout)
	if len(logFileName) > 0 {
		var err error
		logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			log.Warningf("File %s, can't be opened, using STDOUT for logging.", logFileName)
		} else {
			log.SetOutput(logFile)
		}
	}
}

func main() {
	initLogger("")

	if version != "" && buildTime != "" {
		log.Infof("Starting %s v%s build at %s", os.Args[0], version, buildTime)
	}

	webhooks.Start()
	defer webhooks.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
