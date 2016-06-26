package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/pat"
	"github.com/russross/blackfriday"
)

// Post is a blog post
type Post struct {
	Body string    `json:"body"`
	Time time.Time `json:"time"`
}

var db []Post

func main() {
	fmt.Println("Hello, world")

	// Ensure db is an empty slice, not nil
	db = []Post{}

	r := pat.New()

	r.Get("/hello", hello)
	r.Post("/markdown", markdown)

	r.Post("/posts", addPost)
	r.Get("/posts", getPosts)
	r.Delete("/posts/{id}", delPost)

	r.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on localhost:" + port)
	http.ListenAndServe(":"+port, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	body := []byte("Hello, world")
	w.Write(body)
}

func markdown(w http.ResponseWriter, r *http.Request) {
	body := []byte(r.FormValue("body"))
	markdown := blackfriday.MarkdownCommon(body)
	w.Write(markdown)
}

// addPost is responsible for adding a new post
func addPost(w http.ResponseWriter, r *http.Request) {
	// Make the Post variable that will hold our input
	var p Post

	// Decode the request into variable p
	err := json.NewDecoder(r.Body).Decode(&p)

	// If decoding failed, give the user an error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.Time = time.Now()

	db = append(db, p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Print(err)
	}
}

// getPosts lists all posts as JSON
func getPosts(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// delPost removes a post
func delPost(w http.ResponseWriter, r *http.Request) {

	// Figure out which post they want to delete
	idStr := r.URL.Query().Get(":id")

	// Convert their input to an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make sure it's a Post that exists
	if id < 0 || id >= len(db) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	db = append(db[:id], db[id+1:]...)

	w.WriteHeader(http.StatusNoContent)
}
