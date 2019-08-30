package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	bookID, err := strconv.ParseInt(param.ByName("bookID"), 10, 64)
	if err != nil {
		log.Printf("[internal][GetUserById] fail to convert user_id into int :%+v", err)
		renderJSON(w, []byte(`
		{
			"message":"k4m03 n4k4l"
		}
		`), http.StatusBadRequest)
		return
	}
	// TODO: Implement this. Query = SELECT id, title, author, isbn, stock FROM books WHERE id = <bookID>
	query := fmt.Sprintf("SELECT id,title,author,isbn,stock FROM books WHERE id=$1")

	rows, err := h.DB.Query(query, bookID)
	if err != nil {
		log.Println(err)
		return
	}
	var books []Book
	for rows.Next() {
		book := Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Stock)
		if err != nil {
			log.Println(err)
			return
		}
		books = append(books, book)
	}
	bytes, err := json.Marshal(books)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, bytes, http.StatusOK)

}

// InsertBook a function to insert book to DB
func (h *Handler) InsertBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: Implement this. Query = INSERT INTO books (id, title, author, isbn, stock) VALUES (<id>, '<title>', '<author>', '<isbn>', <stock>)
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Println(err)
		return
	}
	// executing insert query
	query := fmt.Sprintf("INSERT INTO books (id,title,author,isbn,stock) VALUES (%d,'%s','%s','%s',%d) ", book.ID, book.Title, book.Author, book.ISBN, book.Stock)
	_, err = h.DB.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Insert Book Successfully"
	}
	`), http.StatusOK)

}

// EditBook a function to change book data in DB, with given params
func (h *Handler) EditBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: Implement this. Query = UPDATE books SET title = '<title>', author = '<author>', isbn = '<isbn>', stock = <stock> WHERE id = <id>
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Println(err)
		return
	}
	// executing insert query
	query := fmt.Sprintf("UPDATE books SET title='%s', author='%s', isbn='%s',stock=%d WHERE id = '%s'", book.Title, book.Author, book.ISBN, book.Stock, param.ByName("bookID"))
	_, err = h.DB.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Book Updated Successfully",
	}
	`), http.StatusOK)
}

// DeleteBookByID a function to remove book data from DB, given bookID
func (h *Handler) DeleteBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = DELETE FROM books WHERE id = <id>
	bookID := param.ByName("bookID")
	query := fmt.Sprintf("DELETE FROM books WHERE id=%s", bookID)
	_, err := h.DB.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Book Deleted Successfully"
	}
	`), http.StatusOK)
}

// InsertMultipleBooks a function to insert multiple book data, given file of books data
func (h *Handler) InsertMultipleBooks(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var buffer bytes.Buffer

	file, header, err := r.FormFile("books")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// get file name
	name := strings.Split(header.Filename, ".")
	if name[1] != "csv" {
		log.Println("File format not supported")
		return
	}
	log.Printf("Received a file with name = %s\n", name[0])

	// copy file to buffer
	io.Copy(&buffer, file)

	// TODO: uncomment this when implementing
	// contents := buffer.String()

	// Split contents to rows
	// TODO: uncomment this when implementing
	// rows := strings.Split(contents, "\n")

	// TODO: iterate csv rows here.

	buffer.Reset()

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert book success!"
	}
	`), http.StatusOK)
}

// LendBook a function to record book lending in DB and update book stock in book tables
func (h *Handler) LendBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this.
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var lr LendRequest
	err = json.Unmarshal(body, &lr)
	if err != nil {
		log.Println(err)
		return
	}

	// Get stock query = SELECT stock FROM books WHERE id = <bookID>
	query := fmt.Sprintf("SELECT stock FROM books WHERE id=%d", lr.BookID)
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	book := Book{}
	for rows.Next() {
		err = rows.Scan(&book.Stock)
		if err != nil {
			log.Println(err)
			return
		}
	}
	// Insert Book Lending query = INSERT INTO lend (user_id, book_id) VALUES (<userID>, <bookID>)
	query2 := fmt.Sprintf("INSERT INTO lend (user_id,book_id) VALUES (%d,%d) ", lr.UserID, lr.BookID)
	_, err = h.DB.Exec(query2)
	if err != nil {
		log.Println(err)
		return
	}
	book.Stock = book.Stock - 1
	// Update stock query = UPDATE books SET stock = <newStock> WHERE id = <bookID>
	query3 := fmt.Sprintf("UPDATE books SET stock = %d WHERE id = %d", book.Stock, lr.BookID)
	_, err = h.DB.Exec(query3)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Stock Book Successfully"
	}
	`), http.StatusOK)
	// Read userID

	// parse json body

	// Get book stock from DB

	// Insert Book to Lend tables

	// Update Book stock query
}
