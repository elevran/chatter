//
package main

import (
	"net/http"
	"os"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mediator"
	app.Flags = Flags
	app.Action = func(context *cli.Context) error {
		config := newConfig(context)
		room := newChatRoom()

		http.HandleFunc("/hello", room.hello)
		http.HandleFunc("/goodbye", room.goodbye)
		http.HandleFunc("/room", room.message)
		return http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPPort), nil)

	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithError(err).Fatalf("Error running main")
	}
}
