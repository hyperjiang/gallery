package main

import (
	"github.com/hyperjiang/gallery-service/app/provider"
	"github.com/hyperjiang/gallery-service/app/router"

	"log"

	"github.com/gin-gonic/gin"
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

	app.StaticFile("/favicon.ico", di.Config().Server.PublicDir+"/favicon.ico")

	router.Route(app)

	// Listen and Serve
	app.Run(":80")
}
