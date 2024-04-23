package helpers

import (
	"time"
)

// Note: Do not change here
// Reference: base/app/orgss/orgs/schema.go
type CustomInput struct {
	ID             int    `gorm:"primaryKey"`
	FormType       string // registration, onboarding ...
	InputName      string
	Required       bool
	InputType      string
	Description    string
	HTMLContent    string
	AnswerFieldMap string
	Options        *string `gorm:"type:json"`
	OrderSort      int
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
