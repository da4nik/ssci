package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"

	"github.com/da4nik/ssci/config"
	"github.com/da4nik/ssci/webhooks"
)

var version string
var buildTime string

var logFile *os.File

// func testRequest() {
// 	data := github.PushEvent{
// 		Pusher: github.User{
// 			Email: "some@aaa.rrr",
// 			Name:  "Some User",
// 		},
// 		Repository: github.Repository{
// 			Name:     "somerepo",
// 			FullName: "da4nik/somerepo",
// 			CloneURL: "git@github.com:da4nik/ssci.git",
// 		},
// 	}
//
// 	json, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	_, err = http.Post("http://webhooks.makstep.ru/github", "application/json", bytes.NewReader(json))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

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

	// Creating worspace dir if not exists
	os.MkdirAll(config.Workspace, 0755)

	webhooks.Start()
	defer webhooks.Stop()

	// time.Sleep(2 * time.Second)
	// testRequest()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
