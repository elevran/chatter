package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mediator"
	app.Flags = Flags
	app.Action = func(context *cli.Context) error {
		logrus.SetLevel(logrus.DebugLevel)

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
