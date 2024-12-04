package api

import (
	"ip/bilibili"
	"ip/config"
	"ip/ip"

	"github.com/gin-gonic/gin"
)

func InitHandleRouter(cfg *config.Config, router *gin.Engine) {
	apiRouter := router.Group("api")
	{
		apiRouter.GET("/healthcheck", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		apiRouter.Any("/ip-lookup", func(c *gin.Context) {
			ip.IPHandler(c)
		})
		apiRouter.Any("/ip", func(c *gin.Context) {
			ip.IPPureHandler(c)
		})
		apiRouter.Any("/bilibili", func(c *gin.Context) {
			bilibili.Bilibili(c)
		})
	}
}
