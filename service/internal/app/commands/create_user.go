package commands

import (
	"context"

	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
)

type CreateUserHandler CommandHandler[CreateUser, *users.User]

type CreateUser struct {
	ID          int64
	Password    string
	IsSuperuser bool
	Username    string
	FirstName   string
	LastName    string
	Email       string
}

type createUserHandler struct {
	repository users.Repository
}

func (c createUserHandler) Handle(ctx context.Context, cmd CreateUser) (*users.User, error) {
	newUser, err := c.repository.Create(
		ctx,
		users.User{
			ID:          uint(cmd.ID),
			Password:    cmd.Password,
			IsSuperuser: cmd.IsSuperuser,
			Username:    cmd.Username,
			FirstName:   cmd.FirstName,
			LastName:    cmd.LastName,
		})

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func NewCreateUserHandler(o users.Repository) CreateUserHandler {
	return createUserHandler{
		repository: o,
	}
}
