package route

import (
	"github.com/gin-gonic/gin"
)

// Router ...
type Router interface {
	Route() []Route
}

// Route ...
type Route struct {
	Name         string
	Method       string
	RelativePath string
	Auth         bool
	HandlerFunc  gin.HandlerFunc
}
