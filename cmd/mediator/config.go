package main

import "github.com/urfave/cli"

var Flags = []cli.Flag{
	cli.StringFlag{
		Name:   "room_id",
		EnvVar: "ROOM_ID",
		Value:  "",
		Usage:  "Game On registration room id",
	},

	cli.StringFlag{
		Name:   "room_service_url",
		EnvVar: "ROOM_SERVICE_URL",
		Value:  "http://localhost:6379/room",
		Usage:  "Room service URL",
	},

	cli.IntFlag{
		Name:   "websocket_port",
		EnvVar: "WEBSOCKET_PORT",
		Value:  3000,
		Usage:  "Port to listen for incoming websocket connections",
	},
}

type Config struct {
	RoomID         string
	RoomServiceURL string
	WebsocketPort  int
}

func newConfig(context *cli.Context) *Config {
	return &Config{
		RoomID:         context.String("room_id"),
		RoomServiceURL: context.String("room_service_url"),
		WebsocketPort:  context.Int("websocket_port"),
	}
}
