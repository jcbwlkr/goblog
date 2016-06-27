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
)

func main() {

	// r is our router from gorilla/pat
	r := pat.New()

	// Map some requests to handlers
	r.Post("/posts", addPost)
	r.Get("/posts", getPosts)
	r.Delete("/posts/{id}", delPost)

	// Serve static content from the public directory directly
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// Figure out the port we should bind to. For local development it
	// will default to 8080 but on Heroku the PORT environment variabl
	// // Figure out the port we should bind to. For local development
	// it will default to 8080 but on Heroku we will use the PORT
	// environment variablee
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start listening for requests
	fmt.Println("Listening on localhost:" + port)
	err := http.ListenAndServe(":"+port, r)

	// If there was a problem binding to the port then err will be set.
	// Log the message and kill the application with an exit code.
	if err != nil {
		log.Fatal(err)
	}
}

// addPost is responsible for adding a new post
func addPost(w http.ResponseWriter, r *http.Request) {
	// Make the Post variable that will hold our input
	var p Post

	// Decode the request into variable p. If decoding failed, give the
	// user an error.
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We do not expect the front end to send us the current time so we
	// set it ourselves.
	p.Time = time.Now()

	// Add the new Post to the database. Calling mutex.Lock before will make
	// this goroutine wait until it can get exclusive access to the db. We must
	// call Unlock when we're done or subsequent calls will wait forever.
	mutex.Lock()
	db = append(db, p)
	mutex.Unlock()

	// Tell the front end we succesfully created the post and respond
	// back with it's full JSON.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Print(err)
	}
}

// getPosts lists all posts as JSON
func getPosts(w http.ResponseWriter, r *http.Request) {

	// We can simply encode the db to our ResponseWriter but we must ensure we
	// have exclusive access first. Because there is more than one return path
	// from this function we will defer the Unlock call so it will happen at
	// the end of this function's execution regardless of where it returns.
	mutex.Lock()
	defer mutex.Unlock()

	if err := json.NewEncoder(w).Encode(db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// delPost removes a post
func delPost(w http.ResponseWriter, r *http.Request) {

	// Figure out which post they want to delete. gorilla/pat puts the URL
	// parameters in the request's URL.Query. We pull that out then convert it
	// to an int.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Before we read or write to db we must get the lock. Defer the Unlock
	// since there are multiple return paths.
	mutex.Lock()
	defer mutex.Unlock()

	// Make sure it's a Post that exists
	if id < 0 || id >= len(db) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	db = append(db[:id], db[id+1:]...)

	w.WriteHeader(http.StatusNoContent)
}
