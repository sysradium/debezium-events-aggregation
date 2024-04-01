package commands

import (
	"context"

	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
)

type DeleteUserHandler CommandHandler[DeleteUser, *users.User]

type DeleteUser struct {
	ID int64
}

type deleteUserHandler struct {
	repository users.Repository
}

func (c deleteUserHandler) Handle(ctx context.Context, cmd DeleteUser) (*users.User, error) {
	if err := c.repository.Delete(ctx, cmd.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

func NewDeleteUserHandler(o users.Repository) deleteUserHandler {
	return deleteUserHandler{
		repository: o,
	}
}
