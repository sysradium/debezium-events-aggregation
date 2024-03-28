package ephemeral

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
)

var _ users.Repository = (*Ephemeral)(nil)

type Ephemeral struct {
	m *sync.RWMutex
	s map[string]users.User

	userFactory users.Factory
}

func (e *Ephemeral) Create(_ context.Context, o users.User) (*users.User, error) {
	e.m.Lock()
	defer e.m.Unlock()

	userID := uuid.New()
	o.ID = userID
	e.s[userID.String()] = o

	return &o, nil

}

func (e *Ephemeral) Get(_ context.Context, id uuid.UUID) (*users.User, error) {
	e.m.RLock()
	defer e.m.RUnlock()

	if o, ok := e.s[id.String()]; ok {
		return &o, nil
	}

	return nil, users.ErrNotFound
}

func (e *Ephemeral) List(_ context.Context) ([]*users.User, error) {
	e.m.RLock()
	defer e.m.RUnlock()

	users := make([]*users.User, 0, len(e.s))

	for _, o := range e.s {
		o := o
		users = append(users, &o)
	}

	return users, nil
}

func (e *Ephemeral) Update(ctx context.Context, id uuid.UUID, updateFn users.UpdaterFn) (*users.User, error) {
	user, err := e.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	e.m.Lock()
	defer e.m.Unlock()

	updateduser, err := updateFn(user)
	if err != nil {
		return nil, err
	}

	e.s[id.String()] = *updateduser

	return updateduser, nil
}

func New() *Ephemeral {
	return &Ephemeral{
		m: &sync.RWMutex{},
		s: map[string]users.User{},
	}
}
