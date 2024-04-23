package seeds

import (
	"base/core/users"
)

func DefaultUsers() []users.User {
	return []users.User{
		{
			FirstName: "Admin",
			LastName:  "Gashi",
			Email:     "admin@base.base.al",
			Password:  "admin",
			Active:    true,
		},
		{
			FirstName: "Reader",
			LastName:  "Gashi",
			Email:     "reader@base.base.al",
			Password:  "reader",
			Active:    true,
		},
	}
}
