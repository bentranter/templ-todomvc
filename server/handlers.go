package server

import (
	"net/http"
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

	sd := getSessionData(r)
	if !sd.ShouldAutofocus && len(sd.Todos) == 0 {
		sd.ShouldAutofocus = true
	}

	remaining := 0
	for _, todo := range sd.Todos {
		if todo.State != "completed" {
			remaining++
		}
	}

	components.
		Page(components.PageProps{
			Todos:               sd.Todos,
			Filter:              filter,
			Remaining:           remaining,
			Completed:           len(sd.Todos) - remaining,
			ShouldAutofocus:     sd.ShouldAutofocus,
			PreserveQueryParams: components.PreserveQueryParams(r),
		}).
		Render(r.Context(), w)
}

// TodoCreateHandler handles the POST request to create new todos.
func TodoCreateHandler(w http.ResponseWriter, r *http.Request) {
	text := strings.TrimSpace(r.FormValue("todo"))

	if text != "" {
		sd := getSessionData(r)
		sd.Todos = append(sd.Todos, components.NewTodo(text))
		sd.ShouldAutofocus = true
		saveSessionData(w, r, sd)
	}

	redirect(w, r, "/")
}

func TodoShowEditHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "none"
	}

	sd := getSessionData(r)
	remaining := 0
	for _, todo := range sd.Todos {
		if todo.State != "completed" {
			remaining++
		}
	}

	components.
		Page(components.PageProps{
			Todos:               sd.Todos,
			EditID:              id,
			Filter:              filter,
			Remaining:           remaining,
			Completed:           len(sd.Todos) - remaining,
			ShouldAutofocus:     false, // The form should never autofocus the create form when editing.
			PreserveQueryParams: components.PreserveQueryParams(r),
		}).
		Render(r.Context(), w)
}

func TodoEditHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	// If updated text is provided, save the update.
	if todoText := r.FormValue("text"); todoText != "" {
		sd := getSessionData(r)

		for i, todo := range sd.Todos {
			if todo.ID == id {
				sd.Todos[i].Text = todoText
			}
		}
		sd.ShouldAutofocus = false
		saveSessionData(w, r, sd)
		redirect(w, r, "/")
		return
	}

	// Otherwise flip the state.
	sd := getSessionData(r)
	for i, todo := range sd.Todos {
		if todo.ID == id {
			if sd.Todos[i].State == "active" {
				sd.Todos[i].State = "completed"
			} else {
				sd.Todos[i].State = "active"
			}
		}
	}
	sd.ShouldAutofocus = false
	saveSessionData(w, r, sd)
	redirect(w, r, "/")
}

func TodoDestroyHandler(w http.ResponseWriter, r *http.Request) {
	id := muxpatterns.PathValue(r, "id")

	sd := getSessionData(r)
	for i, todo := range sd.Todos {
		if todo.ID == id {
			sd.Todos = append(sd.Todos[:i], sd.Todos[i+1:]...)
			break
		}
	}
	sd.ShouldAutofocus = false
	saveSessionData(w, r, sd)
	redirect(w, r, "/")
}

func TodoClearCompletedHandler(w http.ResponseWriter, r *http.Request) {
	sd := getSessionData(r)

	// Find all of the todos that are deleted and therefore must be deleted.
	completedTodoIDs := make([]string, 0)
	for _, todo := range sd.Todos {
		if todo.State == "completed" {
			completedTodoIDs = append(completedTodoIDs, todo.ID)
		}
	}
	// Iterate over the list and delete the completed todos from the list.
	for _, completedTodoID := range completedTodoIDs {
		for i, todo := range sd.Todos {
			if todo.ID == completedTodoID {
				sd.Todos = append(sd.Todos[:i], sd.Todos[i+1:]...)
			}
		}
	}

	sd.ShouldAutofocus = false
	saveSessionData(w, r, sd)
	redirect(w, r, "/")
}

func TodoSelectAllHandler(w http.ResponseWriter, r *http.Request) {
	sd := getSessionData(r)

	todoState := "active"
	for _, todo := range sd.Todos {
		if todo.State == "active" {
			todoState = "completed"
			break
		}
	}

	for i := range sd.Todos {
		sd.Todos[i].State = todoState
	}

	sd.ShouldAutofocus = false
	saveSessionData(w, r, sd)
	redirect(w, r, "/")
}
