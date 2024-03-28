package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	UserID      int64
	Password    string
	LastLogin   *time.Time
	IsSuperuser bool
	Username    string
	FirstName   string
	LastName    string
	Email       string
	IsStaff     bool
	IsActive    bool
	DateJoined  time.Time
}

type Factory struct {
}

func (f Factory) New(userID int64, email string) (*User, error) {
	return &User{
		ID:     uuid.New(),
		UserID: userID,
		Email:  email,
	}, nil
}
