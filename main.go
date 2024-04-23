package main

import "github.com/gin-gonic/gin"

var router *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	err := DatabaseConnection()
	if err != nil {
		return
	}
	router.RedirectTrailingSlash = true
	//redirect "/articles" to "/"
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/articles")

	})
	articleRoutes := router.Group("/articles")
	{
		articleRoutes.GET("/", getArticles)
		articleRoutes.GET("/:id", getArticle)
		articleRoutes.POST("/", addArticle)
		articleRoutes.PUT("/:id", updateArticle)
		articleRoutes.DELETE("/:id", deleteArticle)
	}
	// Start serving the application
	err = router.Run()
	if err != nil {
		return
	}
}
