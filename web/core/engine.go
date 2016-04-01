package core

import (
	"github.com/gin-gonic/gin"
	"github.com/itpkg/husky/web"
	"github.com/jinzhu/gorm"
)

//Engine core engine
type Engine struct {
}

//Mount web mount point
func (p *Engine) Mount(*gin.RouterGroup) {

}

//Migrate database migrate
func (p *Engine) Migrate(*gorm.DB) {

}

func init() {
	web.Register(&Engine{})
}
