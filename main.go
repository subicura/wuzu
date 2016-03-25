package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/subicura/wuzu/build"
	"github.com/subicura/wuzu/config"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "wuzu"
	app.Usage = "build tool with docker"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: fmt.Sprintf("Init %s project", app.Name),
			Action: func(c *cli.Context) {
				config.CreateConfigFile()
			},
		},
		{
			Name:  "build",
			Usage: "Run build",
			Action: func(c *cli.Context) {
				build.Build()
			},
		},
	}
	app.Run(os.Args)
}
