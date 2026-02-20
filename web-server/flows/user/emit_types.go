package user

import "time"

// User is the user record returned by user flows.
type User struct {
	ID        string
	IdpID     string
	UserType  string
	Status    string
	CreatedAt time.Time
}
