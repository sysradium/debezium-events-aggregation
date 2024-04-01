package users

import (
	"context"
)

type UpdaterFn func(*User) (*User, error)

type Repository interface {
	Getter
	Lister
	Creator
	Updater
	Deleter
}

type Getter interface {
	Get(context.Context, int64) (*User, error)
}

type Lister interface {
	List(context.Context) ([]*User, error)
}

type Creator interface {
	Create(context.Context, User) (*User, error)
}

type Updater interface {
	Update(context.Context, int64, UpdaterFn) (*User, error)
}

type Deleter interface {
	Delete(context.Context, int64) error
}
