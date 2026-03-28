package model

import "time"

// Book represents a book in the bookstore.
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	AuthorID    string    `json:"author_id"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Author represents a book author.
type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Review represents a book review.
type Review struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	Reviewer  string    `json:"reviewer"`
	CreatedAt time.Time `json:"created_at"`
}

// BookFilter holds optional filter parameters for listing books.
type BookFilter struct {
	AuthorID string
	Title    string
	MinPrice *float64
	MaxPrice *float64
	SortBy   string // "title", "price", "published_at", "created_at"
	Order    string // "asc", "desc"
	Page     int
	PerPage  int
}

// PagedResult wraps a list response with pagination metadata.
type PagedResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}
