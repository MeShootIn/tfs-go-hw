package main

import (
	"encoding/json"
	"io"
	"lection05/chat"
	"lection05/user"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	CookieLogin    = "LOGIN"
	CookiePassword = "PASSWORD"
	UserLogin      = "login"
	UserPassword   = "password"
	IP             = "localhost"
	PORT           = "5000"
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
	messages.Get("/general", GetMessagesGeneral)
	messages.Post("/general", PostMessagesGeneral)
	messages.Get("/me", GetMessagesMe)
	messages.Post("/me", PostMessagesMe)
	root.Mount("/messages", messages)

	log.Fatal(http.ListenAndServe(":"+PORT, root))
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

	if json.Unmarshal(body, &u) != nil || user.ValidateUser(u) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if userDB.Register(u) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

// curl --cookie "LOGIN=admin;PASSWORD=kek" localhost:5000/messages/general
// curl --cookie "LOGIN=Sasha;PASSWORD=lol" localhost:5000/messages/general
func GetMessagesGeneral(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("*** GENERAL CHAR ***:\n\n" + generalChat.String()))
}

type Msg struct {
	Text string `json:"text"`
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" -H "Content-Type: application/json" --data '{"text":"Hello"}' -X POST localhost:5000/messages/general
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" -H "Content-Type: application/json" --data '{"text":"World"}' -X POST localhost:5000/messages/general
func PostMessagesGeneral(w http.ResponseWriter, r *http.Request) {
	login, err := r.Cookie(CookieLogin)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	msg := Msg{}

	if json.Unmarshal(body, &msg) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	generalChat.SendMessage(login.Value, msg.Text)
}

type MeParams struct {
	With string `json:"with"`
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" -H "Content-Type: application/json" --data '{"with":"Sasha"}' -X POST localhost:5000/messages/me
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" -H "Content-Type: application/json" --data '{"with":"admin"}' -X POST localhost:5000/messages/me
func GetMessagesMe(w http.ResponseWriter, r *http.Request) {
	from, err := r.Cookie(CookieLogin)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	meParams := MeParams{}

	if json.Unmarshal(body, &meParams) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, _ = w.Write([]byte("*** CHAT WITH \"" + meParams.With + "\" ***\n\n" + personalChat[from.Value][meParams.With].String()))
}

type PrsMsg struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

// curl -v --cookie "LOGIN=admin;PASSWORD=kek" -H "Content-Type: application/json" --data '{"text":"Hello!","to":"Sasha"}' -X POST localhost:5000/messages/me
// curl -v --cookie "LOGIN=Sasha;PASSWORD=lol" -H "Content-Type: application/json" --data '{"text":"Hello!","to":"admin"}' -X POST localhost:5000/messages/me
func PostMessagesMe(w http.ResponseWriter, r *http.Request) {
	from, err := r.Cookie(CookieLogin)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	prsMsg := PrsMsg{}

	if json.Unmarshal(body, &prsMsg) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	personalChat.SendMessage(from.Value, prsMsg.To, prsMsg.Text)
}
