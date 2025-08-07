// UserManager handles user authentication and profile management
package main

import (
	"strings"
	"time"
)

type User struct {
    ID       int       `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

func (u *User) ValidateEmail() bool {
    return strings.Contains(u.Email, "@") && strings.Contains(u.Email, ".")
}

func CreateUser(username, email string) *User {
    return &User{
        ID:       generateID(),
        Username: username,
        Email:    email,
        Created:  time.Now(),
    }
}

func generateID() int {
    // Simple ID generation (in real app, use proper UUID)
    return int(time.Now().Unix())
}
