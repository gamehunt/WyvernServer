package http

import (
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
	"wyvern/server/internal/transport/http/handler"
	"wyvern/server/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func addRoutes(router     *gin.Engine, 
			   logger     *slog.Logger, 
			   authSvc    *service.AuthService,
		       sessionSvc *service.SessionService,
		   	   tokenSvc   *service.TokenService) {
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	authHandler := handler.NewAuthHandler(logger, authSvc, sessionSvc)

	public := router.Group("/auth")
    {
        public.POST("/login",    authHandler.Login)
        public.POST("/register", authHandler.Register)
        public.POST("/refresh",  authHandler.Refresh)
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
