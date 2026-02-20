package user

import "time"

// User is the app_users record returned by the repository.
type User struct {
	ID        string
	IdpID     string
	UserType  string
	Status    string
	CreatedAt time.Time
}
