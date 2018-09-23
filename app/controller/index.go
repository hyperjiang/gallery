package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gallery-service/app/model"
	"github.com/hyperjiang/gallery-service/app/provider"
)

// IndexController - the default controller
type IndexController struct{}

// Index - the default page
func (ctrl *IndexController) Index(c *gin.Context) {
	var files model.Files
	err := files.Get()
	if err != nil {
		provider.DI().Log().Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"files": files,
	})
}

// Version - show the version
func (ctrl *IndexController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": provider.DI().Config().Server.Version,
	})
}
