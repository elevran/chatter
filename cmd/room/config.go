package main

import "github.com/urfave/cli"

var Flags = []cli.Flag{
	cli.IntFlag{
		Name:   "http_port",
		EnvVar: "HTTP_PORT",
		Value:  3000,
		Usage:  "Port to listen for incoming websocket connections",
	},
}

type Config struct {
	HTTPPort int
}

func newConfig(context *cli.Context) *Config {
	return &Config{
		HTTPPort: context.Int("http_port"),
	}
}
