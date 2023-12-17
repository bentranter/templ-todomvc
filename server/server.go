package server

import (
	"fmt"
	"net/http"

	"github.com/bentranter/templ-todomvc/components"
	"github.com/jba/muxpatterns"
)

var state = []components.Todo{}

func Start() {
	mux := muxpatterns.NewServeMux()

	mux.HandleFunc("GET /{$}", HomeHandler)
	mux.HandleFunc("POST /todos", TodoCreateHandler)
	mux.HandleFunc("GET /todos/{id}", TodoShowEditHandler)
	mux.HandleFunc("POST /todos/{id}", TodoEditHandler)
	mux.HandleFunc("POST /todos/{id}/destroy", TodoDestroyHandler)
	mux.HandleFunc("POST /todos/clear", TodoClearCompletedHandler)
	mux.HandleFunc("POST /todos/select", TodoSelectAllHandler)

	mux.HandleFunc("GET /css/app.css", RenderFileHandler)
	mux.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}
