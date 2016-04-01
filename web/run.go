package web

import (
	"os"

	"github.com/codegangsta/cli"
)

// Run Main entry
func Run() error {
	app := cli.NewApp()
	app.Name = "husky"
	app.Usage = "IT-PACKAGE web framework"
	app.Version = "v20160329"
	app.Commands = []cli.Command{}
	for _, en := range engines {
		cmd := en.Shell()
		app.Commands = append(app.Commands, cmd...)
	}

	return app.Run(os.Args)
}
