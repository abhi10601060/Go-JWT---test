package main

import (
	"apiwithmysql/auth"
	"apiwithmysql/model"
	"apiwithmysql/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/book/{isbn}", AddBook).Methods("POST")
	r.HandleFunc("/book/{isbn}", getBookByIsbn).Methods("GET")
	r.HandleFunc("/book/{isbn}", updateBook).Methods("PUT")
	r.HandleFunc("/books", getAllBooks).Methods("GET")
	r.HandleFunc("/book/{isbn}", deleteBookbyIsbn).Methods("DELETE")

	r.HandleFunc("/signup", signUp).Methods("POST")
	r.HandleFunc("/authorize", authorize).Methods("GET")
	r.HandleFunc("/login", login).Methods("POST")

	defer log.Fatal(http.ListenAndServe(":9000", r))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	auth.Signup(w, r)
}

func authorize(w http.ResponseWriter, r *http.Request) {
	auth.Authorise(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	auth.Login(w, r)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	isbn := params["isbn"]

	var book model.Book

	json.NewDecoder(r.Body).Decode(&book)
	book.Isbn = isbn

	//Add to the db
	res := repo.Addbook(&book)
	if res != 1 {
		log.Fatal("book not saved successfully..")
	}
	json.NewEncoder(w).Encode(&book)
	fmt.Println("Book saved successfully...")
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isAuthorised := auth.Authorise(w, r)
	if !isAuthorised {
		json.NewEncoder(w).Encode("You are not authorised user...")
		return
	}

	var books []model.Book

	repo.GetAllbooks(&books)

	json.NewEncoder(w).Encode(&books)
}

func deleteBookbyIsbn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isAuthorised := auth.Authorise(w, r)
	if !isAuthorised {
		json.NewEncoder(w).Encode("You are not authorised user...")
		return
	}

	params := mux.Vars(r)

	isbn := params["isbn"]
	fmt.Println("delete book called with isbn : ", isbn)
	repo.DeleteBookByIsbn(isbn)

	json.NewEncoder(w).Encode("book deleted successfully..")
}

func getBookByIsbn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isAuthorised := auth.Authorise(w, r)
	if !isAuthorised {
		json.NewEncoder(w).Encode("You are not authorised user...")
		return
	}

	params := mux.Vars(r)
	isbn := params["isbn"]

	var book model.Book
	res, err := repo.GetBookByIsbn(isbn, &book)

	if res == 1 {
		json.NewEncoder(w).Encode(&book)
	} else {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(err.Error())
	}

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isAuthorised := auth.Authorise(w, r)
	if !isAuthorised {
		json.NewEncoder(w).Encode("You are not authorised user...")
		return
	}

	params := mux.Vars(r)
	isbn := params["isbn"]

	var book model.Book
	repo.GetBookByIsbn(isbn, &book)
	json.NewDecoder(r.Body).Decode(&book)
	repo.UpdateBook(&book)

	json.NewEncoder(w).Encode(&book)
}
