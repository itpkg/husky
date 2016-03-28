package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Engine web engine
type Engine interface {
	Mount(*gin.RouterGroup)
	Migrate(*gorm.DB)
}
