package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gallery-service/app/provider"
	"github.com/hyperjiang/gallery-service/app/router"

	"log"
)

func main() {

	// create a dependency injection container
	di := provider.DI()

	// check if the configs can be loaded successfully
	if di.Config() == nil {
		log.Fatal("Fail to load configs")
	}

	// init db
	if di.InitDB() != nil {
		log.Fatal("Fail to init database")
	}

	app := gin.Default()

	app.LoadHTMLGlob(di.Config().Server.ViewDir + "/*")
	app.StaticFile("/favicon.ico", di.Config().Server.PublicDir+"/favicon.ico")

	router.Route(app)

	// Listen and Serve
	app.Run(":80")
}
