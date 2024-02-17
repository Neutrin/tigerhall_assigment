package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nitin/tigerhall/core/controllers"
	"github.com/nitin/tigerhall/core/inits"
	"github.com/nitin/tigerhall/core/repositiories"
	"github.com/nitin/tigerhall/core/routes"
	"github.com/nitin/tigerhall/core/services"
)

func main() {
	userRepo, err := repositiories.NewPostgresqlUserRepo()
	if err != nil {
		panic(err)
	}
	tigerRepo, err := repositiories.NewPostgresqlTigerRepo()
	if err != nil {
		panic(err)
	}
	inits.InitStorageClient()
	log.Println(" intilaisation succesfull")
	userNotifyService := services.NewSightingNotification(2, userRepo, services.NewMailNotification())
	InitRoutes(controllers.NewAuthController(userRepo),
		controllers.NewTigerController(tigerRepo, userRepo, userNotifyService))

}

func InitRoutes(authController controllers.AuthController, testingController controllers.TigerControllers) {
	r := gin.Default()
	routes.Routes(r, authController, testingController)
	//routes.Routes(r, testingController)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
