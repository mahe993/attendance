/*
Package services provides business logic for performing requests specific to each endpoint.
*/
package services

import (
	"fmt"
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

type AuthService struct{}

var (
	Auth AuthService
)

// Variables that holds the auth states
var (
	MapUsers    = map[string]User{}
	MapSessions = map[string]string{}
)

func init() {
	// init special access for admin
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	MapUsers["admin"] = User{"admin", bPassword, "admin", "admin"}
}

func (*AuthService) Login(res http.ResponseWriter, req *http.Request) {
	// process form submission
	loginID := req.FormValue("loginID")
	password := req.FormValue("password")

	// check if user exist with loginID
	myUser, ok := MapUsers[loginID]

	if !ok {
		http.Error(res, "Login ID and/or password do not match", http.StatusUnauthorized)
		return
	}

	// Matching of password entered
	err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
	if err != nil {
		http.Error(res, "Login ID and/or password do not match", http.StatusForbidden)
		return
	}

	// create session cookie
	id := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:  "myCookie",
		Value: id.String(),
		Path:  "/",
	}

	http.SetCookie(res, myCookie)
	MapSessions[myCookie.Value] = loginID

	http.Redirect(res, req, fmt.Sprintf("/?user=%s", myUser.ID), http.StatusSeeOther)
}

func (*AuthService) Logout(res http.ResponseWriter, req *http.Request) {
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(MapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func GetUser(res http.ResponseWriter, req *http.Request) User {
	// if the user exists already, get user
	var myUser User

	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Println(err)
		}
		return myUser
	}

	if loginID, ok := MapSessions[myCookie.Value]; ok {
		myUser = MapUsers[loginID]
	}

	return myUser
}
