package user

import (
	"errors"
	"strings"
)

type Login = string
type Password = string

type User struct {
	Login    *Login    `json:"login"`
	Password *Password `json:"password"`
}

var (
	ErrNoLogin           = errors.New("login was not passed")
	ErrNoPassword        = errors.New("password was not passed")
	ErrAlreadyRegistered = errors.New("a user with this login is already registered")
	ErrWrongLogin        = errors.New("wrong login")
	ErrWrongPassword     = errors.New("wrong password")
)

func ValidateLogin(login *Login) error {
	if login == nil {
		return ErrNoLogin
	}

	if strings.TrimSpace(*login) == "" {
		return ErrWrongLogin
	}

	return nil
}

func ValidatePassword(password *Password) error {
	if password == nil {
		return ErrNoPassword
	}

	if strings.TrimSpace(*password) == "" {
		return ErrWrongPassword
	}

	return nil
}

func ValidateUser(user User) error {
	if err := ValidateLogin(user.Login); err != nil {
		return err
	}

	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

type DataBase map[Login]User

func (db DataBase) Register(user User) error {
	if err := ValidateUser(user); err != nil {
		return err
	}

	if _, ok := db[*user.Login]; ok {
		return ErrAlreadyRegistered
	}

	db[*user.Login] = User{
		Login:    user.Login,
		Password: user.Password,
	}

	return nil
}

func (db DataBase) Login(user User) error {
	if err := ValidateUser(user); err != nil {
		return err
	}

	u, ok := db[*user.Login]

	if !ok {
		return ErrWrongLogin
	}

	if *user.Password != *u.Password {
		return ErrWrongPassword
	}

	return nil
}

func (db DataBase) GetUser(login string) *User {
	u, ok := db[login]

	if !ok {
		return nil
	}

	return &u
}
