package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World! Welcome to Books CRUD API")
	})
	r.POST("api/books", createBook)
	r.GET("api/books", getAllBooks)
	r.GET("api/books/:id", getBookByID)
	r.PUT("api/books/:id", updateBookByID)
	r.DELETE("api/books/:id", deleteBookByID)

	r.Run(":80")
}

func createTable() {
	query := `
        CREATE TABLE IF NOT EXISTS books (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            author TEXT NOT NULL
        )
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", "The Great Gatsby", "A novel by F. Scott Fitzgerald")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", "To Kill a Mockingbird", "A novel by Harper Lee")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", "1984", "A novel by George Orwell")
	if err != nil {
		log.Fatal(err)
	}
}

func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", book.Title, book.Author)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := result.LastInsertId()
	book.ID = int(id)

	c.JSON(http.StatusCreated, book)
}

func getAllBooks(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func getBookByID(c *gin.Context) {
	id := getIDParam(c)

	var book Book
	err := db.QueryRow("SELECT * FROM books WHERE id=?", id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			log.Fatal(err)
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

func updateBookByID(c *gin.Context) {
	id := getIDParam(c)

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE books SET title=?, author=? WHERE id=?", book.Title, book.Author, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			log.Fatal(err)
		}
		return
	}

	book.ID = id
	c.JSON(http.StatusOK, book)
}

func deleteBookByID(c *gin.Context) {
	id := getIDParam(c)
	_, err := db.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			log.Fatal(err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func getIDParam(c *gin.Context) int {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
