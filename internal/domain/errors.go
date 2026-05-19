package domain

import "errors"

var (
	//ErrInvalidBookData is the error returned when the book data is invalid.
	ErrInvalidBookData = errors.New("invalid book data")
	//ErrBookNotFound is the error returned when a book is not found.
	ErrBookNotFound = errors.New("book not found")
	//ErrUnauthorized is the error returned when the user is not authorized to perform an action.
	ErrUnauthorized = errors.New("unauthorized")
	//ErrInternalServer is the error returned when an internal server error occurs.
	ErrInternalServer = errors.New("internal server error")
)