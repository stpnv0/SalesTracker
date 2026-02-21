package router

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type itemHandler interface {
	Create(c *ginext.Context)
	List(c *ginext.Context)
	Update(c *ginext.Context)
	Delete(c *ginext.Context)
	GetByID(c *ginext.Context)
}

type analyticsHandler interface {
	Get(c *ginext.Context)
}

type exportHandler interface {
	CSV(c *ginext.Context)
}

func InitRouter(
	mode string,
	itemHandler itemHandler,
	analyticsHandler analyticsHandler,
	exportHandler exportHandler,
	mw ...ginext.HandlerFunc,
) *ginext.Engine {
	router := ginext.New(mode)
	router.Use(ginext.Recovery())
	router.Use(mw...)

	api := router.Group("/api")
	{
		api.POST("/items", itemHandler.Create)
		api.GET("/items", itemHandler.List)
		api.GET("/items/:id", itemHandler.GetByID)
		api.PUT("/items/:id", itemHandler.Update)
		api.DELETE("/items/:id", itemHandler.Delete)

		api.GET("/analytics", analyticsHandler.Get)

		api.GET("/export/csv", exportHandler.CSV)
	}

	router.GET("/health", func(c *ginext.Context) {
		c.JSON(http.StatusOK, ginext.H{"status": "ok"})
	})

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "web/static")

	router.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}
