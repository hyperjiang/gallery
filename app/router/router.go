package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gallery-service/app/controller"
)

// Route makes the routing
func Route(app *gin.Engine) {
	indexController := new(controller.IndexController)

	app.GET(
		"/", indexController.Index,
	).GET(
		"/version", indexController.Version,
	)

	fileController := new(controller.FileController)
	app.GET(
		"/upload", fileController.Form,
	).POST(
		"/upload", fileController.Upload,
	).GET(
		"/file", fileController.Read,
	)

}
