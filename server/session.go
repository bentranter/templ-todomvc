package server

import (
	"crypto/rand"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/bentranter/templ-todomvc/components"
	"github.com/gorilla/sessions"
)

type sessionData struct {
	Todos           []components.Todo
	ShouldAutofocus bool
}

var store = func() *sessions.CookieStore {
	key, ok := os.LookupEnv("SESSION_KEY")
	if !ok {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			panic(err)
		}
		key = string(buf)
	}
	gob.Register(&sessionData{})
	return sessions.NewCookieStore([]byte(key))
}()

func getSessionData(r *http.Request) *sessionData {
	session, err := store.Get(r, "_session")
	if err != nil {
		println("failed to get session:", err.Error())
		session = sessions.NewSession(store, "_session")
	}

	v := session.Values["sessiondata"]
	if v == nil {
		println("sessiondata is nil")
		return &sessionData{Todos: make([]components.Todo, 0)}
	}

	sd, ok := v.(*sessionData)
	if !ok {
		println("v is not sessiondata, but is", v)
		return &sessionData{Todos: make([]components.Todo, 0)}
	}
	return sd
}

func saveSessionData(w http.ResponseWriter, r *http.Request, sd *sessionData) {
	session, err := store.Get(r, "_session")
	if err != nil {
		session = sessions.NewSession(store, "_session")
	}

	session.Values["sessiondata"] = sd
	if err := session.Save(r, w); err != nil {
		println("failed to save session:", err.Error())
	}
}
