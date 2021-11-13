package main

import (
	"context"
	"encoding/json"
	"io"
	"lection05/chat"
	"lection05/user"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ContextKey string

const (
	CookieLogin               = "LOGIN"
	CookiePassword            = "PASSWORD"
	UserLogin      ContextKey = "login"
	UserPassword   ContextKey = "password"
	HOST                      = "localhost"
	PORT                      = "5000"
)

var (
	userDB       = user.DataBase{}
	generalChat  = chat.Chat{}
	personalChat = chat.PersonalChat{}
)

func main() {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Post("/login", Login)
	root.Post("/register", Register)

	messages := chi.NewRouter()
	messages.Use(Auth)
	messages.Get("/general", GetMessagesGeneral)
	messages.Post("/general", PostMessagesGeneral)
	messages.Get("/to/{login}", GetMessagesTo)
	messages.Post("/to/{login}", PostMessagesTo)
	root.Mount("/messages", messages)

	log.Fatal(http.ListenAndServe(HOST+":"+PORT, root))
}

func Auth(handler http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		cookieErrStatus := func(cookieErr error) int {
			switch cookieErr {
			case nil:
			case http.ErrNoCookie:
				return http.StatusUnauthorized
			default:
				return http.StatusInternalServerError
			}

			return http.StatusOK
		}

		cookieLogin, cookieLoginErr := r.Cookie(CookieLogin)

		if err := cookieErrStatus(cookieLoginErr); err != http.StatusOK {
			w.WriteHeader(err)

			return
		}

		cookiePassword, cookiePasswordErr := r.Cookie(CookiePassword)

		if err := cookieErrStatus(cookiePasswordErr); err != http.StatusOK {
			w.WriteHeader(err)

			return
		}

		u := user.User{
			Login:    &cookieLogin.Value,
			Password: &cookiePassword.Value,
		}

		if err := userDB.Login(u); err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), UserLogin, *u.Login)
		ctx = context.WithValue(ctx, UserPassword, *u.Password)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}

// curl -v -H "Content-Type: application/json" --data '{"login":"admin","password":"kek"}' -X POST localhost:5000/login
// curl -v -H "Content-Type: application/json" --data '{"login":"Sasha","password":"lol"}' -X POST localhost:5000/login
func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()
	u := user.User{}

	if json.Unmarshal(body, &u) != nil || userDB.Login(u) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	cookieLogin := &http.Cookie{
		Name:  CookieLogin,
		Value: *u.Login,
		Path:  "/",
	}
	http.SetCookie(w, cookieLogin)

	cookiePassword := &http.Cookie{
		Name:  CookiePassword,
		Value: *u.Password,
		Path:  "/",
	}
	http.SetCookie(w, cookiePassword)
}

// curl -v -H "Content-Type: application/json" --data '{"login":"admin","password":"kek"}' -X POST localhost:5000/register
// curl -v -H "Content-Type: application/json" --data '{"login":"Sasha","password":"lol"}' -X POST localhost:5000/register
func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()
	u := user.User{}

	if json.Unmarshal(body, &u) != nil || userDB.Register(u) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	cookieLogin := &http.Cookie{
		Name:  CookieLogin,
		Value: *u.Login,
		Path:  "/",
	}
	http.SetCookie(w, cookieLogin)

	cookiePassword := &http.Cookie{
		Name:  CookiePassword,
		Value: *u.Password,
		Path:  "/",
	}
	http.SetCookie(w, cookiePassword)
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" localhost:5000/messages/general
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" localhost:5000/messages/general
func GetMessagesGeneral(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("*** GENERAL CHAT ***\n\n" + generalChat.String() + "\n*** GENERAL CHAT ***\n"))
}

type Message struct {
	Text string `json:"text"`
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" -H "Content-Type: application/json" --data '{"text":"Hello"}' -X POST localhost:5000/messages/general
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" -H "Content-Type: application/json" --data '{"text":"World"}' -X POST localhost:5000/messages/general
func PostMessagesGeneral(w http.ResponseWriter, r *http.Request) {
	me, ok := r.Context().Value(UserLogin).(user.Login)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()
	message := Message{}

	if json.Unmarshal(body, &message) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	generalChat.SendMessage(me, message.Text)
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" localhost:5000/messages/to/Sasha
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" localhost:5000/messages/to/admin
func GetMessagesTo(w http.ResponseWriter, r *http.Request) {
	me, ok := r.Context().Value(UserLogin).(user.Login)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	to := chi.URLParam(r, "login")

	if userDB.GetUser(to) == nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, _ = w.Write([]byte("*** PERSONAL CHAT WITH \"" + to + "\" ***\n\n" + personalChat[me][to].String() + "\n*** PERSONAL CHAT WITH \"" + to + "\" ***\n"))
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" -H "Content-Type: application/json" --data '{"text":"Hello, Sasha!"}' -X POST localhost:5000/messages/to/Sasha
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" -H "Content-Type: application/json" --data '{"text":"Hello, admin!"}' -X POST localhost:5000/messages/to/admin
func PostMessagesTo(w http.ResponseWriter, r *http.Request) {
	me, ok := r.Context().Value(UserLogin).(user.Login)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()
	message := Message{}

	if json.Unmarshal(body, &message) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	to := chi.URLParam(r, "login")

	if userDB.GetUser(to) == nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	personalChat.SendMessage(me, to, message.Text)
}
