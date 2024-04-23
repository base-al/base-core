package users

import (
	"time"
)

// SeedUsers seeds user data
func Seed() []User {
	return []User{
		{
			Email:         "admin@base.base.al",
			Username:      nil,
			Password:      "admin",
			FirstName:     "Admin",
			LastName:      "Gashi",
			Active:        true,
			Phone:         "1234567890",
			VerifiedEmail: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     nil,
			DeletedAt:     nil,
		},
		{
			Email:         "reader@base.base.al",
			Username:      nil,
			Password:      "reader",
			FirstName:     "Reader",
			LastName:      "Gashi",
			Active:        true,
			Phone:         "9876543210",
			VerifiedEmail: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     nil,
			DeletedAt:     nil,
		},
		// Add more users as needed
	}
}
