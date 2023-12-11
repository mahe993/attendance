/*
Package users provides functionalities for performing requests specific to users.

It includes functionalities for:
  - Login/Logout
  - User struct
  - All methods relevant to users
*/
package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents the metadata of a user
type User struct {
	ID       string
	Password []byte
	First    string
	Last     string
}

// Variables that holds the user states
var (
	mapUsers    = map[string]User{}
	mapSessions = map[string]string{}
)

func init() {
	// init special access for admin
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = User{"admin", bPassword, "admin", "admin"}
}

func GetUser() {
	fmt.Println("USERS")
}
