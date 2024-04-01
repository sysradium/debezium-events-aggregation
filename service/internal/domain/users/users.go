package users

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primarykey"`
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
		ID:    uint(userID),
		Email: email,
	}, nil
}
