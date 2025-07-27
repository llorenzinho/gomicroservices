package server

import "github.com/gin-gonic/gin"

func CreateUserRouter(server *gin.Engine) *gin.RouterGroup {
	api := server.Group("/api")
	users := api.Group("/users")
	return users
}

func CreateJwtRouter(server *gin.Engine) *gin.RouterGroup {
	api := server.Group("/api")
	jwt := api.Group("/jwt")
	return jwt
}
