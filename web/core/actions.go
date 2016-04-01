package core

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

//ENV environment var
var ENV = cli.StringFlag{
	Name:   "environment, e",
	Value:  "development",
	Usage:  "Specifies the environment to run this server under (test/development/production).",
	EnvVar: "ITPKG_ENV",
}

//EnvAction env command
func EnvAction(fn func(string, *cli.Context) error) func(*cli.Context) {
	return func(ctx *cli.Context) {
		log.Println("Begin...")
		if err := fn(ctx.String("environment"), ctx); err == nil {
			log.Println("Done!!!")
		} else {
			log.Fatalln(err)
		}
	}

}

//ConfigAction config command
func ConfigAction(fn func(*Config, *cli.Context) error) func(*cli.Context) {
	return EnvAction(func(env string, ctx *cli.Context) error {
		var cfg Config
		if err := Load(fmt.Sprintf("%s.toml", env), &cfg); err != nil {
			return err
		}
		cfg.Env = env
		return fn(&cfg, ctx)
	})
}
