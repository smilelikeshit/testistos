package domain

import (
	"context"
	"errors"
)

type Animal struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type AnimalRepository interface {
	Store(ctx context.Context, an *Animal) (err error)
	GetByID(ctx context.Context, id int) (Animal, error)
}

type AnimalUseCase interface {
	Store(ctx context.Context, an *Animal) (err error)
	GetByID(ctx context.Context, id int) (res *Animal, err error)
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
)
