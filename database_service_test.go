package main

import "testing"

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

func TestCreateArticleTable(t *testing.T) {
	err := createArticleTable(DB)
	if err != nil {
		t.Fatalf("Failed to create article table: %v", err)
	}
}

func TestInsertArticle(t *testing.T) {
	a := Article{Title: "Test Title", Content: "Test Content"}
	ok, err := insertArticle(a)
	if err != nil {
		t.Fatalf("Failed to insert article: %v", err)
	}

	if !ok {
		t.Fatal("Failed to insert article")
	}
}

func TestGetAllArticles(t *testing.T) {
	articles, err := getAllArticles()
	if err != nil {
		t.Fatalf("Failed to get articles: %v", err)
	}

	if len(articles) == 0 {
		t.Fatal("No articles found")
	}
}

func TestGetArticleByID(t *testing.T) {
	article, err := getArticleByID(1)
	if err != nil {
		t.Fatalf("Failed to get article by ID: %v", err)
	}

	if article.ID != 1 {
		t.Fatalf("Article ID does not match: got %d, want %d", article.ID, 1)
	}
}

func TestUpdateArticle(t *testing.T) {
	// Prepare the updated article
	a := Article{ID: 5, Title: "Updated Title", Content: "Updated Content"}

	// Update the article
	ok, err := updateArticleByID(a, 5)
	if err != nil {
		t.Fatalf("Failed to update article: %v", err)
	}
	if !ok {
		t.Fatal("Update article returned false")
	}
}

func TestDeleteArticle(t *testing.T) {
	// Delete the article
	ok, err := deleteArticleByID(5)
	if err != nil {
		t.Fatalf("Failed to delete article: %v", err)
	}
	if !ok {
		t.Fatal("Delete article returned false")
	}
}
