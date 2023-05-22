package api

import (
	"github.com/casbin/casbin/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.com/micro/api_gateway/api/docs"
	v1 "gitlab.com/micro/api_gateway/api/handlers/v1"
	"gitlab.com/micro/api_gateway/api/middleware"
	"gitlab.com/micro/api_gateway/api/token"
	"gitlab.com/micro/api_gateway/config"
	"gitlab.com/micro/api_gateway/pkg/logger"
	"gitlab.com/micro/api_gateway/services"
	"gitlab.com/micro/api_gateway/storage/repo"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf            config.Config
	Logger          logger.Logger
	ServiceManager  services.IServiceManager
	InMemoryStorage repo.RedisRepo
	CasbinEnforcer  *casbin.Enforcer
}

// New ...
// @title           Mind-Blow
// @version         2.0
// @description     Some description
// @termsOfService  Golang

// @contact.name   Amirkhan
// @telegram    https://t.me/torexanovich
// @contact.email  torexanovich.l@gmail.com

// @host    	   localhost:5050

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	// corConfig := cors.DefaultConfig()
	// corConfig.AllowAllOrigins = true
	// corConfig.AllowCredentials = true
	// corConfig.AllowHeaders = []string{"*"}
	// corConfig.AllowBrowserExtensions = true
	// corConfig.AllowMethods = []string{"*"}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := token.JWTHandler{
		SigninKey: option.Conf.SignInKey,
		Log:       option.Logger,
	}

	router.Use(middleware.NewAuth(option.CasbinEnforcer, jwtHandler, option.Conf))

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.InMemoryStorage,
		JWTHandler:     jwtHandler,
		Casbin: *option.CasbinEnforcer,
	})

	api := router.Group("/v1")

	// users
	// api.POST("/user", handlerV1.CreateUser)
	api.GET("/user", handlerV1.GetUser)
	api.GET("/users", handlerV1.GetAllUsers)
	api.PUT("/user/:id", handlerV1.UpdateUser)
	api.DELETE("/user/:id", handlerV1.DeleteUser)

	// posts
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/:id", handlerV1.GetPost)
	api.GET("/posts/:id", handlerV1.GetPostByUserId)
	api.PUT("/post/:id", handlerV1.UpdatePost)
	api.DELETE("/post/:id", handlerV1.DeletePost)

	// comments
	api.POST("comment", handlerV1.CreateComment)
	api.GET("comment/:id", handlerV1.GetCommentsForPost)
	api.DELETE("comment/:id", handlerV1.DeleteComment)

	// register
	api.POST("/user/register", handlerV1.RegisterUser)
	api.GET("/verify/:email/:code", handlerV1.VerifyUser)
	api.GET("/login/:email/:password", handlerV1.Login)

	// admin
	api.POST("/admin/register", handlerV1.RegisterAdmin)
	api.GET("/admin/verify/{email}/{code}", handlerV1.VerifyAdmin)
	api.PATCH("/admin/create_mod/{id}", handlerV1.CreateMod)

	// swagger
	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
