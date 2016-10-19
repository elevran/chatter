package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	app := cli.NewApp()
	app.Name = "mediator"
	app.Flags = Flags
	app.Action = func(context *cli.Context) error {
		config := newConfig(context)
		s := newServer(config)

		logrus.Infof("Starting mediator service on port %d", config.WebsocketPort)

		http.HandleFunc("/", s.handleHTTP)
		return http.ListenAndServe(fmt.Sprintf(":%d", config.WebsocketPort), nil)

	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithError(err).Fatalf("Error running main")
	}
}
