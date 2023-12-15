/*
Package services provides business logic for performing requests specific to each endpoint.
*/
package services

import (
	"net/http"

	"attendance.com/src/logger"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents the metadata of an authenticated user
type User struct {
	ID       string
	Password []byte
	First    string
	Last     string
}

type AuthService struct {
	currUser *User
}

var (
	Auth AuthService = AuthService{}
)

// Variables that holds all user and session details
var (
	MapUsers    = map[string]*User{}
	MapSessions = map[string]string{}
)

func init() {
	// init special access for admin
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	MapUsers["admin"] = &User{"admin", bPassword, "admin", "admin"}
}

func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	// process form submission
	loginID := r.FormValue("loginID")
	password := r.FormValue("password")

	// check if user exist with loginID
	myUser, ok := MapUsers[loginID]
	if !ok {
		http.Error(w, "Login ID and/or password do not match", http.StatusUnauthorized)
		return
	}

	// Matching of password entered
	err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
	if err != nil {
		http.Error(w, "Login ID and/or password do not match", http.StatusForbidden)
		return
	}

	// create session cookie
	c := createSessCookie(w, r)

	a.currUser = myUser

	// map cookie value to loginID
	MapSessions[<-c] = loginID

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	sessCookie, _ := r.Cookie("sessCookie")
	// delete the session
	delete(MapSessions, sessCookie.Value)
	// remove the cookie
	sessCookie = &http.Cookie{
		Name:   "sessCookie",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, sessCookie)

	a.currUser = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthService) GetUser(r *http.Request) *User {
	// get current session cookie
	sessCookie, err := r.Cookie("sessCookie")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Println(err)
		}
		return a.currUser
	}

	if loginID, ok := MapSessions[sessCookie.Value]; ok {
		a.currUser = MapUsers[loginID]
	}

	return a.currUser
}

func createSessCookie(w http.ResponseWriter, r *http.Request) <-chan string {
	c := make(chan string)

	go func(c chan<- string) {
		id := uuid.NewV4()
		sessCookie := &http.Cookie{
			Name:  "sessCookie",
			Value: id.String(),
			Path:  "/",
		}

		http.SetCookie(w, sessCookie)
		c <- sessCookie.Value
	}(c)

	return c
}
