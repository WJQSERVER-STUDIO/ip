package api

import (
	"ip/bapi"
	"ip/config"
	"ip/ip"

	"github.com/infinite-iroha/touka"
)

func InitHandleRouter(cfg *config.Config, router *touka.Engine) {
	apiRouter := router.Group("api")
	{
		apiRouter.GET("/healthcheck", func(c *touka.Context) {
			c.JSON(200, touka.H{
				"status": "ok",
			})
		})
		apiRouter.GET("/ip-lookup", func(c *touka.Context) {
			ip.IPHandler(c)
		})
		apiRouter.GET("/ip", func(c *touka.Context) {
			ip.IPPureHandler(c)
		})
		apiRouter.GET("/bilibili", func(c *touka.Context) {
			bapi.Bilibili(c)
		})
	}
}
