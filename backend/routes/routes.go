package routes

import (
	"desafio-fullstack/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const account = "/account/"

func AccountHandler() {
	r := gin.Default()
	r.GET(account+"person/:id", controllers.GetPerson)
	r.POST(account+"person", controllers.PostPerson)
	r.PUT(account+"person/:id", controllers.PutPerson)
	r.PATCH(account+"person/friend/:id", controllers.AddFriend)
	r.GET(account+"friends/:id", controllers.FindFriends)
	r.POST(account+"card/:id", controllers.PostCard)
	r.GET(account+"cards/:id", controllers.GetCards)
	r.POST(account+"transfer", controllers.PostTransfer)
	r.GET(account+"bank-statement/:id", controllers.GetTransfers)
	r.GET(account+"bank-statement", controllers.GetTransfersFriends)
	// docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
