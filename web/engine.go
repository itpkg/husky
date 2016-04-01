package web

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Engine web engine
type Engine interface {
	Mount(*gin.RouterGroup)
	Migrate(*gorm.DB)
	Shell() []cli.Command
}

var engines []Engine

//Loop loop in engines
func Loop(fn func(Engine) error) error {
	for _, en := range engines {
		if err := fn(en); err != nil {
			return err
		}
	}
	return nil
}

//Register register engine
func Register(ens ...Engine) {
	engines = append(engines, ens...)
}
