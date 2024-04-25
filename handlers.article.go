package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getArticles(c *gin.Context) {
	articles, err := getAllArticles()
	if err != nil {
		err = c.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			return
		}
	}
	render(c, gin.H{
		"title":   "Home Page",
		"payload": articles,
	}, "index.html")
}

func showCreateArticlePage(c *gin.Context) {
	render(c, gin.H{
		"title": "Create New Article",
	}, "create-article.html")

}

func showEditArticlePage(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("id")); err == nil {
		if article, err := getArticleByID(articleID); err == nil {
			render(c, gin.H{
				"title":   "Edit Article",
				"payload": article,
			}, "edit-article.html")
		} else {
			err = c.AbortWithError(http.StatusNotFound, err)
			if err != nil {
				return
			}
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func getArticle(c *gin.Context) {
	//Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if article, err := getArticleByID(articleID); err == nil {
			// Call the HTML method of the context to render a template
			render(c, gin.H{
				"title":   article.Title,
				"payload": article,
			}, "article.html")
		} else {
			// if the article is not found, abort with an error
			err = c.AbortWithError(http.StatusNotFound, err)
			if err != nil {
				return
			}
		}
	} else {
		// if an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func addArticle(c *gin.Context) {
	if c.Request.Header.Get("Accept") == "application/json" {
		var newArticle Article
		if err := c.BindJSON(&newArticle); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		_, err := insertArticle(newArticle)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully"})
		}
	} else {
		// Get the input from the user
		title := c.PostForm("title")
		content := c.PostForm("content")

		// Validate input
		if title == "" || content == "" {
			err := c.AbortWithError(http.StatusBadRequest, nil)
			if err != nil {
				return
			}
			return
		}

		// Insert data into the database
		newArticle := Article{Title: title, Content: content}
		_, err := insertArticle(newArticle)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/")
		}
	}
}

func updateArticle(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("id")); err == nil {
		if c.Request.Header.Get("Accept") == "application/json" {
			var newArticle Article
			if err := c.BindJSON(&newArticle); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			if _, err := updateArticleByID(newArticle, articleID); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
			}
		}
	}
}

func deleteArticle(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("id")); err == nil {
		if c.Request.Header.Get("Accept") == "application/json" {
			if _, err := deleteArticleByID(articleID); err != nil {
				err := c.AbortWithError(http.StatusInternalServerError, err)
				if err != nil {
					panic(err)
				}
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
			}
		}
	}
}

// Render one of HTML, JSON, CSV based on the 'Accept' header of the request
// if the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		//Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		//Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
