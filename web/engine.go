package web

import (
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Mount(*gin.RouterGroup)
}
