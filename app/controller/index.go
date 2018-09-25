package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Unknwon/paginater"
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gallery-service/app/model"
	"github.com/hyperjiang/gallery-service/app/provider"
)

// IndexController - the default controller
type IndexController struct{}

// Index - the default page
func (ctrl *IndexController) Index(c *gin.Context) {

	t := c.Query("t")
	if t == "" {
		t = "image"
	}

	// get the total count of files by type
	total, err := model.CountFilesByType(t)
	if err != nil {
		provider.DI().Log().Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(total)

	limit := uint(100)

	// calculate the db offset
	var page uint64 = 1
	if c.Query("p") != "" {
		page, err = strconv.ParseUint(c.Query("p"), 10, 64)
		if err != nil {
			page = 1
		}
	}
	offset := (uint(page) - 1) * limit

	var numPages = 5
	// Arguments:
	// - Total number of rows
	// - Number of rows in one page
	// - Current page number
	// - Number of page links to be displayed
	p := paginater.New(int(total), int(limit), int(page), numPages)

	var files model.Files
	err = files.GetByType(t, limit, offset)
	if err != nil {
		provider.DI().Log().Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"files": files,
		"type":  t,
		"page":  p,
		"n1":    numPages - 2,
		"n2":    p.TotalPages() - 2,
	})
}

// Version - show the version
func (ctrl *IndexController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": provider.DI().Config().Server.Version,
	})
}
