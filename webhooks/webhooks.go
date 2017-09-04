package webhooks

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/ssci/types/github"
)

type server struct {
	mux *http.ServeMux
}

var httpServer *http.Server

// Start starts listening webhooks
func Start() {
	srv := &server{mux: http.NewServeMux()}
	srv.mux.HandleFunc("/github", srv.github)

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":2020"
	}

	httpServer = &http.Server{Addr: addr, Handler: srv}
	log().Infof("Listening on http://0.0.0.0%s", httpServer.Addr)

	httpServer.ListenAndServe()
}

// Stop stops listening webhooks
func Stop() {
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log().Infof("Shutdown with timeout: %s\n", timeout)

	if err := httpServer.Shutdown(ctx); err != nil {
		log().Errorf("Error: %v", err)
	} else {
		log().Info("Server stopped")
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Server", "example Go server")
	s.mux.ServeHTTP(w, r)
}

func (s server) github(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	decoder := json.NewDecoder(req.Body)
	var githubPushNotification github.PushEvent
	if err := decoder.Decode(&githubPushNotification); err != nil {
		log().Errorf("Error decoding github notification: %v", err)
		return
	}

	// TODO: #8 Process github notifications
	// spew.Dump(githubPushNotification)
}

func log() *logrus.Entry {
	return logrus.WithField("module", "webhooks")
}
