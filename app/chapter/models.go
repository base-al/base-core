package chapter

import "time"

// Chapter represents a single CMS chapter with comprehensive details including metadata and status.
type Chapter struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`    // e.g., "draft", "published", "archived"
	AuthorID    int       `json:"author_id"` // ID of the user who created or last updated the chapter
}

// CreateChapterRequest model for the create API call, now includes AuthorID.
type CreateChapterRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
}

// CreateChapterResponse model for the create response.
type CreateChapterResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// ReadChapterRequest model for fetching a chapter.
type ReadChapterRequest struct {
	ID int `json:"id"`
}

// ReadChapterResponse model for returning a chapter.
type ReadChapterResponse struct {
	Chapter
}

// UpdateChapterRequest model for updating chapter details.
type UpdateChapterRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`
	AuthorID    int       `json:"author_id"`
}

// UpdateChapterResponse model for the update response.
type UpdateChapterResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// DeleteChapterRequest model for deleting a chapter.
type DeleteChapterRequest struct {
	ID int `json:"id"`
}

// DeleteChapterResponse model for the delete response.
type DeleteChapterResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// ListChapterRequest model for listing chapters, can include filters for status.
type ListChapterRequest struct {
	Chapter     int    `json:"chapter"`
	Limit    int    `json:"limit"`
	Status   string `json:"status"`
	AuthorID int    `json:"author_id"`
}

// ListChapterResponse model for the list response.
type ListChapterResponse struct {
	Chapters []Chapter `json:"chapters"`
}
