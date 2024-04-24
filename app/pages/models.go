package pages

import "time"

// Page represents a single CMS page with comprehensive details including metadata and status.
type Page struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`    // e.g., "draft", "published", "archived"
	AuthorID    int       `json:"author_id"` // ID of the user who created or last updated the page
}

// CreatePageRequest model for the create API call, now includes AuthorID.
type CreatePageRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
}

// CreatePageResponse model for the create response.
type CreatePageResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// ReadPageRequest model for fetching a page.
type ReadPageRequest struct {
	ID int `json:"id"`
}

// ReadPageResponse model for returning a page.
type ReadPageResponse struct {
	Page
}

// UpdatePageRequest model for updating page details.
type UpdatePageRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`
	AuthorID    int       `json:"author_id"`
}

// UpdatePageResponse model for the update response.
type UpdatePageResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// DeletePageRequest model for deleting a page.
type DeletePageRequest struct {
	ID int `json:"id"`
}

// DeletePageResponse model for the delete response.
type DeletePageResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// ListPageRequest model for listing pages, can include filters for status.
type ListPageRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Status   string `json:"status"`
	AuthorID int    `json:"author_id"`
}

// ListPageResponse model for the list response.
type ListPageResponse struct {
	Pages []Page `json:"pages"`
}
