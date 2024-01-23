package storage

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	Email    string `bson:"email"`
	HashPass []byte `bson:"hash_pass"`
}
