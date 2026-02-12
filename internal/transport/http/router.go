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
		protected.GET("/test", func(c *gin.Context) {
			uid, _ := c.Get("userId")
			sid, _ := c.Get("sessionId")
			c.JSON(200, gin.H{"UserId": uid, "SessionId": sid})
		})
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
