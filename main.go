package main

import (
	"github.com/urfave/cli/v2"
	cmd "lazyman/com"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Name:  "YseNode~",
		Usage: "My CLI Application",
		Commands: []*cli.Command{
			cmd.Node(),
			cmd.GetALlAlert(),
			cmd.UpdsteALlAlert(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
