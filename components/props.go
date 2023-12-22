package components

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/a-h/templ"
)

//go:generate templ generate

type PageProps struct {
	Todos               []Todo
	EditID              string
	Filter              string
	Remaining           int
	Completed           int
	ShouldAutofocus     bool
	PreserveQueryParams func(s string) templ.SafeURL
}

func PreserveQueryParams(r *http.Request) func(s string) templ.SafeURL {
	return func(s string) templ.SafeURL {
		return templ.URL(s + "?" + r.URL.RawQuery)
	}
}

type Todo struct {
	ID    string // Randomly generated unique ID.
	Text  string // Body of the todo.
	State string // Either "active" or "completed".
}

func NewTodo(text string) Todo {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	id := "todoid_" + hex.EncodeToString(b)

	return Todo{
		ID:    id,
		Text:  text,
		State: "active",
	}
}
