package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var sqlDatabase = "sqlite3"
var sqlDatabasePath = "./articles.db"

// DatabaseConnection is the function that will be called to establish a connection to the database
func DatabaseConnection() error {
	// Open the database
	db, err := sql.Open(sqlDatabase, sqlDatabasePath)
	if err != nil {
		fmt.Printf("Failed to open the database: %v\n", err)
		return err
	}
	// Create the article table
	if err := createArticleTable(db); err != nil {
		fmt.Printf("Failed to create articles table: %v\n", err)
		err := db.Close()
		if err != nil {
			fmt.Printf("Failed to close the database: %v\n", err)
			return err
		}
		// Close the database if there is an error
		return err
	}
	DB = db

	return nil
}

// createArticleTable is a function that creates the articles table in the database if it does not exist
func createArticleTable(db *sql.DB) error {
	// Create the articles table if not exists
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS articles (
	  	id INTEGER PRIMARY KEY AUTOINCREMENT,
	  	title TEXT,
	  	content TEXT
		)`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	statement, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}

	statement.Exec()

	return nil
}

// insertArticle is a function that inserts an article into the database
func insertArticle(newArticle Article) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO articles(title, content) VALUES(?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newArticle.Title, newArticle.Content)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil

}

// getAllArticles is a function that gets all articles from the database
func getAllArticles() ([]Article, error) {
	// Get all articles from the database
	rows, err := DB.Query("SELECT id, title, content FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil
}

// getArticleByID is a function that gets an article by ID from the database
func getArticleByID(id int) (Article, error) {
	// Get an article by ID from the database
	sqlstmt, err := DB.Prepare("SELECT id, title, content FROM articles WHERE id = ?")

	if err != nil {
		return Article{}, err
	}
	a := Article{}

	sqlErr := sqlstmt.QueryRow(id).Scan(&a.ID, &a.Title, &a.Content)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Article{}, nil
		}
		return Article{}, sqlErr
	}
	return a, nil
}

// updateArticle is a function that updates an article in the database
func updateArticleByID(article Article, id int) (bool, error) {
	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE articles SET title = ?, content = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(article.Title, article.Content, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil

}

// deleteArticle is a function that deletes an article from the database
func deleteArticleByID(id int) (bool, error) {
	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE FROM articles WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
