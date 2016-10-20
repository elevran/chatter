package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

func main() {
	logrus.Infof("Starting profaned service")

	finder, err := newProfanityFinder()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to create handler")
	}

	h := &handler{
		profanity: finder,
	}

	http.HandleFunc("/", h.handle)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		logrus.WithError(err).Fatalf("Error running main")
	}
}
