package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/russross/blackfriday"
)

func main() {
	fmt.Println("Hello, world")

	r := pat.New()

	r.Get("/hello", hello)
	r.Post("/markdown", markdown)

	r.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on localhost:" + port)
	http.ListenAndServe(":"+port, r) // Changed "nil" to "r"
}

func hello(w http.ResponseWriter, r *http.Request) {
	now := []byte(time.Now().Format(time.StampMicro))
	if err := ioutil.WriteFile("time", now, os.ModeAppend); err != nil {
		http.Error(w, "could not write file", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	body := []byte("Hello, world")
	w.Write(body)
}

func markdown(w http.ResponseWriter, r *http.Request) {
	body := []byte(r.FormValue("body"))
	markdown := blackfriday.MarkdownCommon(body)
	w.Write(markdown)
}
