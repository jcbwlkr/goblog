package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/russross/blackfriday"
)

// Post is a blog post.
type Post struct {
	Title string    `json:"title"`
	Body  Markdown  `json:"body"`
	Time  time.Time `json:"time"`
}

// Markdown extends string with extra Marshalling behavior.
type Markdown string

// MarshalJSON will turn a Markdown string into HTML for representation in an
// API.
func (m Markdown) MarshalJSON() ([]byte, error) {
	mkd := blackfriday.MarkdownCommon([]byte(m))

	js, err := json.Marshal(string(mkd))
	if err != nil {
		return nil, err
	}

	return js, nil
}

// db is our fake database. It is just a slice of Posts. We want to
// declare it at the package level so it's available throughout this
// package. At this scope we have to declare variables with var which
// will initialize it to the type's zero value.
var db []Post

// mutex is a guard to ensure that only one goroutine can access db at a time.
// If two goroutines accessed db at the same time and one was performing a
// write we would have a race condition which would lead to corrupt or
// unexpected behavior.
var mutex sync.Mutex

// init runs before main if it exists. We do this to ensure db is an
// empty slice, not nil (the zero value for slices).
func init() {
	db = []Post{}
}
