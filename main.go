package main

import "github.com/gin-gonic/gin"

var router *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Load the templates defined in the templates folder
	router.LoadHTMLGlob("templates/*")
	router.RedirectTrailingSlash = true

	// Initialize database connection
	err := DatabaseConnection()
	if err != nil {
		return
	}

	//redirect "/articles" to "/"
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/articles")

	})

	// Define the route for the index page and display the index.html template
	articleRoutes := router.Group("/articles")
	{
		articleRoutes.GET("/", getArticles)
		articleRoutes.GET("/:id", getArticle)
		articleRoutes.GET("/new", showCreateArticlePage)
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
