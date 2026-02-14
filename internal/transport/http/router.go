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

// @title           Wyvern API
// @version         1.0
// @description     Wyvern messaging server API

// @host      localhost:5000
// @BasePath  /
func addRoutes(router     *gin.Engine, 
			   logger     *slog.Logger, 
			   authSvc    *service.AuthService,
		       sessionSvc *service.SessionService,
			   tokenSvc   *service.TokenService,

			   guildsSvc  *service.GuildService,
			   membersSvc *service.MemberService,
		) {


	authHandler   := handler.NewAuthHandler(logger, authSvc, sessionSvc)
	guildsHandler := handler.NewGuildHandler(logger, guildsSvc, membersSvc)

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
		guilds := protected.Group("guilds")
		{
			guilds.POST("/", guildsHandler.CreateGuild)
		}

		guild := guilds.Group(":guildId")	
		{
			guild.GET("/", guildsHandler.GetGuild)
		}

		members := guild.Group("members")
		{
			members.GET("/", func(c *gin.Context) {
				c.String(200, "members " + c.Param("guildId"))
			})

			members.GET("/:userId", func(c *gin.Context) {
				c.String(200, "member: " + c.Param("userId"))
			})
		}

		channels := guild.Group("channels")
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
	guildsSvc  *service.GuildService,
	membersSvc *service.MemberService,
) http.Handler {
	r := gin.Default()

	addRoutes(r, logger, authSvc, 
	sessionSvc, tokenSvc, guildsSvc, membersSvc)

	return r
}
