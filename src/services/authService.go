/*
Package services provides business logic for performing requests specific to each endpoint.
*/
package services

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationPageVariables struct {
	User states.User
	Tab  string
}

type AuthService struct {
	currUser  states.User
	Variables RegistrationPageVariables
}

type Attendance struct {
	LoginID     string
	CheckInTime time.Time
}

var (
	Auth AuthService = AuthService{}
)

func init() {
	// init special access for admin
	logger.Println("Initializing admin user")
	bPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.MinCost)
	states.MapUsers["admin"] = states.User{
		ID:       "admin",
		Password: bPassword,
		First:    "admin",
		Last:     "admin",
	}
	logger.Println("Success!")
}

func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	// process form submission
	loginID := r.FormValue("loginID")
	password := r.FormValue("password")

	// check if user exist with loginID
	myUser, ok := states.MapUsers[loginID]
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
	states.MapSessions[<-c] = loginID

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	sessCookie, _ := r.Cookie("sessCookie")
	// delete the session
	delete(states.MapSessions, sessCookie.Value)
	// remove the cookie
	sessCookie = &http.Cookie{
		Name:   "sessCookie",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, sessCookie)

	a.currUser = states.User{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthService) GetUser(r *http.Request) states.User {
	// get current session cookie
	sessCookie, err := r.Cookie("sessCookie")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Println(err)
		}
		return a.currUser
	}

	if loginID, ok := states.MapSessions[sessCookie.Value]; ok {
		a.currUser = states.MapUsers[loginID]
	}

	return a.currUser
}

func (a *AuthService) RegisterPage(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	a.Variables.User = currUser
	a.Variables.Tab = strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]

	err := templates.Tpl.ExecuteTemplate(w, "registrationPage", a.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

func (a *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("error updating users.json::" + err.(error).Error())
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}()

	// process form submission
	loginID := r.FormValue("loginID")
	password := r.FormValue("password")

	// check if user exist with loginID
	logger.Println(states.MapUsers)
	user, ok := states.MapUsers[loginID]
	logger.Println(user)

	if ok {
		// check if user already has a password
		if len(user.Password) > 0 {
			http.Error(w, "Login ID already registered, try signing in instead.", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Login ID not recognized.", http.StatusUnauthorized)
		return
	}

	// hash the password
	bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Password hashing error:", err)
		return
	}

	// register user
	user.Password = bPassword
	states.MapUsers[loginID] = user

	// Write MapUsers state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.MapUsers, "users.json")
	if err != nil {
		logger.Println(err)
	}

	http.Redirect(w, r, "/auth/success", http.StatusSeeOther)
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
