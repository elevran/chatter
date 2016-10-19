package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "mediator"
	app.Flags = Flags
	app.Action = func(context *cli.Context) error {
		config := newConfig(context)
		s := newServer(config)

		http.HandleFunc("/", s.handleHTTP)
		return http.ListenAndServe(fmt.Sprintf(":%d", config.WebsocketPort), nil)

	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("failure running main: %s", err.Error())
	}

}
