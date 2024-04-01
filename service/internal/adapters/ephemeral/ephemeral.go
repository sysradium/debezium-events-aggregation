package ephemeral

import (
	"context"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ users.Repository = (*Ephemeral)(nil)

type Ephemeral struct {
	userFactory users.Factory
	db          *gorm.DB
}

func (e *Ephemeral) Create(_ context.Context, o users.User) (*users.User, error) {
	if err := e.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&o).Error; err != nil {
		return nil, err
	}

	return &o, nil

}

func (e *Ephemeral) Get(_ context.Context, id int64) (*users.User, error) {
	var user users.User
	if err := e.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (e *Ephemeral) List(_ context.Context) ([]*users.User, error) {
	var users []*users.User
	if err := e.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (e *Ephemeral) Update(ctx context.Context, id int64, updateFn users.UpdaterFn) (*users.User, error) {
	user, err := e.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	updateduser, err := updateFn(user)
	if err != nil {
		return nil, err
	}

	if err := e.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return updateduser, nil
}

func (e *Ephemeral) Delete(_ context.Context, id int64) error {
	if err := e.db.Delete(&users.User{}, id).Error; err != nil {
		return err
	}

	return nil

}
func New() (*Ephemeral, error) {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	// TODO: fix this
	if err := db.AutoMigrate(&users.User{}); err != nil {
		return nil, fmt.Errorf("automigration for user failed: %w", err)
	}
	return &Ephemeral{
		db: db,
	}, nil
}
