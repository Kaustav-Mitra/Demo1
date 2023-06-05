package main

import (
	"database/sql"
	"encoding/json"
	"strconv"

	// "fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var db *sql.DB

// var books []Book

// func GetBooks(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(books)
// }

// func GetBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for _, book := range books {
// 		if book.ID == params["id"] {
// 			json.NewEncoder(w).Encode(book)
// 			return
// 		}
// 	}
// 	json.NewEncoder(w).Encode(nil)
// }

// func CreateBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var book Book
// 	_ = json.NewDecoder(r.Body).Decode(&book)
// 	books = append(books, book)
// 	json.NewEncoder(w).Encode(book)
// }

// func UpdateBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Cotent-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, book := range books {
// 		if book.ID == params["id"] {
// 			books = append(books[:index], books[index+1:]...)
// 			var updatedBook Book
// 			_ = json.NewDecoder(r.Body).Decode(&updatedBook)
// 			updatedBook.ID = params["id"]
// 			books = append(books, updatedBook)
// 			json.NewEncoder(w).Encode(updatedBook)
// 			return
// 		}
// 	}
// 	json.NewEncoder(w).Encode(nil)
// }

// func DeleteBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, book := range books {
// 		if book.ID == params["id"] {
// 			books = append(books[:index], books[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(books)
// }

func main() {

	db, err := sql.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/mydatabase")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected To the database")

	// books = append(books, Book{ID: 1, Title: "Test Book 1", Author: "Test Author 1", Price: 200})
	// books = append(books, Book{ID: 2, Title: "Test Book 2", Author: "Test Author 2", Price: 400})

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooksHandler).Methods("GET")
	router.HandleFunc("/books/{id}", getBookHandler).Methods("GET")
	router.HandleFunc("/books", createBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}", updateBookHandler).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBookHandler).Methods("DELETE")

	// it is to start the http server
	log.Println("Server Endpoint Hit and Sarted localhost:8080 use http://localhost:8080 for testing")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := getBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, books, http.StatusOK)
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	book, err := getBook(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.NotFound(w, r)
		return
	}
	jsonResponse(w, book, http.StatusOK)
}

func createBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = insertBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, book, http.StatusCreated)
}
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = bookID
	err = updateBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, book, http.StatusOK)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	err = deleteBook(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getBooks() ([]Book, error) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func getBook(bookID int) (*Book, error) {

	stmt, err := db.Prepare("SELECT * FROM books WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var book Book
	err = stmt.QueryRow(bookID).Scan(&book.ID, &book.Title, &book.Author, &book.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &book, nil
}

func insertBook(book Book) error {

	stmt, err := db.Prepare("INSERT INTO books (title, author, price) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Title, book.Author, book.Price)
	if err != nil {
		return err
	}

	return nil
}

func updateBook(book Book) error {

	stmt, err := db.Prepare("UPDATE books SET title = ?, author = ?, price = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Title, book.Author, book.Price, book.ID)
	if err != nil {
		return err
	}

	return nil
}

func deleteBook(bookID int) error {

	stmt, err := db.Prepare("DELETE FROM books WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookID)
	if err != nil {
		return err
	}

	return nil
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
