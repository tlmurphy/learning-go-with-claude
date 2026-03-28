package repository

import (
	"context"
	"errors"

	"learning-go-with-claude/projects/03-bookstore-api/internal/model"
)

// Common errors returned by repositories.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// BookRepository defines the data access interface for books.
type BookRepository interface {
	FindAll(ctx context.Context, filter model.BookFilter) ([]model.Book, int, error)
	FindByID(ctx context.Context, id string) (model.Book, error)
	Create(ctx context.Context, book model.Book) (model.Book, error)
	Update(ctx context.Context, id string, book model.Book) (model.Book, error)
	Delete(ctx context.Context, id string) error
}

// AuthorRepository defines the data access interface for authors.
type AuthorRepository interface {
	FindAll(ctx context.Context) ([]model.Author, error)
	FindByID(ctx context.Context, id string) (model.Author, error)
	Create(ctx context.Context, author model.Author) (model.Author, error)
	Update(ctx context.Context, id string, author model.Author) (model.Author, error)
	Delete(ctx context.Context, id string) error
}

// ReviewRepository defines the data access interface for reviews.
type ReviewRepository interface {
	FindByBookID(ctx context.Context, bookID string) ([]model.Review, error)
	FindByID(ctx context.Context, id string) (model.Review, error)
	Create(ctx context.Context, review model.Review) (model.Review, error)
	Update(ctx context.Context, id string, review model.Review) (model.Review, error)
	Delete(ctx context.Context, id string) error
}
