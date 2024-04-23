package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConnection(t *testing.T) {
	err := DatabaseConnection()
	// add some fake data to the database

	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	if DB == nil {
		t.Fatal("DB is nil")
	}

	if err != nil {
		t.Fatalf("Failed to insert article: %v", err)
	}
}

// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an unauthenticated user
func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Ensure the router has the same setup for /articles as the main application
	articleRoutes := r.Group("/articles")
	{
		articleRoutes.GET("/", getArticles)
	}

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/articles/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Home Page"
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), "<title>Home Page</title>") && strings.Contains(string(p), "<h2>Test Article</h2>")

		return statusOK && pageOK
	})
}

func TestArticleUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/articles/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Test Article</title>") > 0

		return statusOK && pageOK
	})
}

func TestArticleListJSON(t *testing.T) {
	r := getRouter(false)

	articleRoutes := r.Group("/articles")
	{
		articleRoutes.GET("/", getArticles)
	}
	req, _ := http.NewRequest("GET", "/articles/", nil)
	req.Header.Add("Accept", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		return statusOK
	})
}

func TestArticleListXML(t *testing.T) {
	r := getRouter(false)

	articleRoutes := r.Group("/articles")
	{
		articleRoutes.GET("/", getArticles)
	}

	req, _ := http.NewRequest("GET", "/articles/", nil)
	req.Header.Add("Accept", "application/xml")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		return statusOK
	})
}

func TestArticleUnauthenticatedJSON(t *testing.T) {
	r := getRouter(true)

	r.GET("/articles/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	req.Header.Add("Accept", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		return statusOK
	})
}

func TestArticleUnauthenticatedXML(t *testing.T) {
	r := getRouter(false)

	r.GET("/articles/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	req.Header.Add("Accept", "application/xml")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		return statusOK
	})
}
