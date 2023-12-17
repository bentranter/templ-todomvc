package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bentranter/templ-todomvc/components"
	"github.com/jba/muxpatterns"
)

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	code := http.StatusFound

	if r.Method == http.MethodPost ||
		r.Method == http.MethodPut ||
		r.Method == http.MethodPatch ||
		r.Method == http.MethodDelete {
		code = http.StatusSeeOther
	}

	if r.URL.RawQuery != "" {
		url = url + "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, url, code)
}

// HomeHandler renders the home page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "none"
	}

	remaining := 0
	for _, todo := range state {
		if todo.State != "completed" {
			remaining++
		}
	}

	components.
		Page(components.PageProps{
			Todos:               state,
			Filter:              filter,
			Remaining:           remaining,
			Completed:           len(state) - remaining,
			PreserveQueryParams: components.PreserveQueryParams(r),
		}).
		Render(r.Context(), w)
}

// TodoCreateHandler handles the POST request to create new todos.
func TodoCreateHandler(w http.ResponseWriter, r *http.Request) {
	todo := strings.TrimSpace(r.FormValue("todo"))

	if todo != "" {
		state = append(state, components.NewTodo(todo))
	}

	redirect(w, r, "/")
}

func TodoShowEditHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "none"
	}

	remaining := 0
	for _, todo := range state {
		if todo.State != "completed" {
			remaining++
		}
	}

	components.
		Page(components.PageProps{
			Todos:               state,
			EditID:              id,
			Filter:              filter,
			Remaining:           remaining,
			Completed:           len(state) - remaining,
			PreserveQueryParams: components.PreserveQueryParams(r),
		}).
		Render(r.Context(), w)
}

func TodoEditHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	// If updated text is provided, save the update.
	if todoText := r.FormValue("text"); todoText != "" {
		for i, todo := range state {
			if todo.ID == id {
				state[i].Text = todoText
			}
		}
		redirect(w, r, "/")
		return
	}

	// Otherwise flip the state.
	for i, todo := range state {
		if todo.ID == id {
			if state[i].State == "active" {
				state[i].State = "completed"
			} else {
				state[i].State = "active"
			}
		}
	}

	redirect(w, r, "/")
}

func TodoDestroyHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	for i, todo := range state {
		if todo.ID == id {
			state = append(state[:i], state[i+1:]...)
			break
		}
	}
	redirect(w, r, "/")
}

func TodoClearCompletedHandler(w http.ResponseWriter, r *http.Request) {
	for i, todo := range state {
		if todo.State == "completed" {
			state = append(state[:i], state[i+1:]...)
		}
	}
	redirect(w, r, "/")
}

func TodoSelectAllHandler(w http.ResponseWriter, r *http.Request) {
	todoState := "active"

	for _, todo := range state {
		if todo.State == "active" {
			todoState = "completed"
			break
		}
	}

	for i := range state {
		state[i].State = todoState
	}

	redirect(w, r, "/")
}

// RenderFileHandler attempts to render the file with the path matching that
// of the incoming request.
func RenderFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := "./" + strings.TrimPrefix(r.URL.Path, "/")

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[error] failed to open %s: %v\n", filename, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/css")
	if _, err := io.Copy(w, f); err != nil {
		log.Printf("[error] failed to write %s: %v\n", filename, err)
	}
}
