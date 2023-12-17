/*
Package services provides business logic for performing requests specific to each endpoint.

Additionally, the package initializes an admin user during the package's initialization.
*/
package services

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegistrationPageVariables is a struct that represents the variables that are passed to the registration page template
type RegistrationPageVariables struct {
	User states.User
	Tab  string
}

// AuthService is a struct that provides methods for handling business logics for requests to the /auth endpoint
type AuthService struct {
	Variables   RegistrationPageVariables
	VariablesMu sync.Mutex
}

// Attendance is a struct that represents a user's attendance record
type Attendance struct {
	LoginID     string
	CheckInTime time.Time
}

// Auth is a global variable that provides access to the AuthService methods
var (
	Auth AuthService = AuthService{}
)

func init() {
	// init special access for admin
	logger.Println("Initializing admin user")
	bPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.MinCost)
	states.SetMapUser("admin", states.User{
		ID:       "admin",
		Password: bPassword,
		First:    "admin",
		Last:     "admin",
	})
	logger.Println("Success!")
}

// The Login method handles the processing of form submissions for user login.
// It checks the provided login ID and password, compares the password hash, and creates a session cookie upon successful login.
// If the login is unsuccessful, the user is redirected to the login page with an error message.
func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	// process form submission
	loginID := r.FormValue("loginID")
	password := r.FormValue("password")

	// check if user exist with loginID
	myUser, ok := states.GetMapUser(loginID)
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

	// map cookie value to loginID
	states.SetMapSession(<-c, loginID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles the processing of user logout requests.
// It deletes the session cookie and redirects the user to the login page.
func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("sessCookie")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Println(err)
		}
	} else {
		// delete the session
		// this map operation is thread safe because cookie values are unique
		delete(states.MapSessions, sessCookie.Value)
	}

	// remove the cookie
	sessCookie = &http.Cookie{
		Name:   "sessCookie",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, sessCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetUser returns the user associated with the current session cookie.
// If no session cookie is found, an empty user is returned.
func (a *AuthService) GetUser(r *http.Request) states.User {
	user := states.User{}
	// get current session cookie
	sessCookie, err := r.Cookie("sessCookie")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Println(err)
		}
		return user
	}

	if loginID, ok := states.GetMapSession(sessCookie.Value); ok {
		if usr, ok := states.GetMapUser(loginID); ok {
			user = usr
		}
	}

	return user
}

// RegisterPage renders the registration page template with the appropriate variables.
func (a *AuthService) RegisterPage(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)

	// Mutex lock to ensure thread-safe access to shared Variables field
	a.VariablesMu.Lock()
	a.Variables.User = currUser
	a.Variables.Tab = strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]

	err := templates.Tpl.ExecuteTemplate(w, "registrationPage", a.Variables)
	a.VariablesMu.Unlock()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

// Register handles the processing of form submissions for user registration.
// It checks if the user already exists, and if not, it hashes the password and registers the user.
// If the user already exists, the user is redirected to the login page with an error message.
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
	user, ok := states.GetMapUser(loginID)
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
	states.SetMapUser(loginID, user)

	// Write MapUsers state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.GetAllMapUsers(), "users.json")
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
