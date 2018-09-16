package controller

import (
	"net/http"

	"github.com/hyperjiang/gallery-service/app/provider"

	"github.com/gin-gonic/gin"
)

// IndexController is the default controller
type IndexController struct{}

// Index the default page
func (ctrl *IndexController) Index(c *gin.Context) {

	provider.DI().Log().Infof("hello %s", "world")

	c.JSON(http.StatusOK, gin.H{
		"version": provider.DI().Config().Server.Version,
	})
}
