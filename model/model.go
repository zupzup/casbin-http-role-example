package model

import "errors"

// User is a user
type User struct {
	ID   int
	Name string
	Role string
}

// Users is a list of users
type Users []User

// Exists checks if a user with the given id exists in the list
func (u Users) Exists(id int) bool {
	exists := false
	for _, user := range u {
		if user.ID == id {
			return true
		}
	}
	return exists
}

// FindByName returns the user with the given name, or returns an error
func (u Users) FindByName(name string) (User, error) {
	for _, user := range u {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, errors.New("USER_NOT_FOUND")
}
