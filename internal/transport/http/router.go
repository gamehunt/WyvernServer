package http

import (
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
	"wyvern/server/internal/transport/http/handler"
	"wyvern/server/internal/transport/http/middleware"

	_ "wyvern/server/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// HealthCheck godoc
// @Summary      Health check
// @Description  Checks if server running
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  
// @Failure      500
// @Router       /health [get]
func addRoutes(router     *gin.Engine, 
			   logger     *slog.Logger, 
			   authSvc    *service.AuthService,
		       sessionSvc *service.SessionService,
		   	   tokenSvc   *service.TokenService) {


	authHandler := handler.NewAuthHandler(logger, authSvc, sessionSvc)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	auth := router.Group("/auth")
    {
        auth.POST("/login",    authHandler.Login)
        auth.POST("/register", authHandler.Register)
        auth.POST("/refresh",  authHandler.Refresh)
    }

    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware(tokenSvc))
	{
		guilds := protected.Group("guilds/:guildId")	
		{
			guilds.GET("/", func(c *gin.Context) {
				c.String(200, c.Param("guildId"))
			})
		}

		members := guilds.Group("members")
		{
			members.GET("/", func(c *gin.Context) {
				c.String(200, "members " + c.Param("guildId"))
			})

			members.GET("/:userId", func(c *gin.Context) {
				c.String(200, "member: " + c.Param("userId"))
			})
		}

		channels := guilds.Group("channels")
		{
			channels.GET("/", func(c *gin.Context) {
				c.String(200, "channels " + c.Param("guildId"))
			})
		}


		channel := channels.Group(":channelId")
		{
			channel.GET("/", func(c *gin.Context) {
				c.String(200, "channel: " + c.Param("channelId"))
			})

		}

		messages := channel.Group("messages")
		{
			messages.GET("/", func(c *gin.Context) {
				c.String(200, "messages: " + c.Param("channelId"))
			})

			messages.GET("/:messageId", func(c *gin.Context) {
				c.String(200, "message: " + c.Param("messageId"))
			})
		}
	}
}

func NewRouter(
	logger     *slog.Logger,
	authSvc    *service.AuthService,
	sessionSvc *service.SessionService,
	tokenSvc   *service.TokenService,
) http.Handler {
	r := gin.Default()

	addRoutes(r, logger, authSvc, sessionSvc, tokenSvc)

	return r
}
