package ginx

import (
	"avocado.com/internal/dao"
	"avocado.com/internal/handler"
	"avocado.com/internal/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, d *dao.Dao, c model.URLConfig) {
	sessionAuth := NewSessionAuth(d, c)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{c.App}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("google_oauth_redirect", handler.HandleGoogleRedirect)
			user.GET("settings", sessionAuth.Auth, handler.HandleGetUserSettings)
		}

		slack := api.Group("/slack")
		{
			slack.POST("events", handler.HandleEventRequest)
			slack.GET("redirect", handler.HandleRedirectRequest)
		}

		linear := api.Group("/linear")
		{
			linear.POST("webhook", handler.HandleLinearWebhook)
			linear.GET("redirect", handler.HandleLinearRedirectRequest)
		}
	}
}
