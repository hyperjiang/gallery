package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gallery-service/app/controller"
	"github.com/hyperjiang/gallery-service/app/provider"
)

// Route makes the routing
func Route(app *gin.Engine) {
	indexController := new(controller.IndexController)

	app.GET(
		"/version", indexController.Version,
	)

	fileController := new(controller.FileController)
	app.GET(
		"/file", fileController.Read,
	)

	accounts := provider.DI().Config().Admin.Accounts()
	authorized := app.Group("/", gin.BasicAuth(accounts))
	authorized.GET(
		"", indexController.Index,
	).GET(
		"upload", fileController.Form,
	).POST(
		"upload", fileController.Upload,
	)
}
