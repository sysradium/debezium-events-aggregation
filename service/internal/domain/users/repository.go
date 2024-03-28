package users

import (
	"context"

	"github.com/google/uuid"
)

type UpdaterFn func(*User) (*User, error)

type Repository interface {
	Getter
	Lister
	Creator
	Updater
}

type Getter interface {
	Get(context.Context, uuid.UUID) (*User, error)
}

type Lister interface {
	List(context.Context) ([]*User, error)
}

type Creator interface {
	Create(context.Context, User) (*User, error)
}

type Updater interface {
	Update(context.Context, uuid.UUID, UpdaterFn) (*User, error)
}
