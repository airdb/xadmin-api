package data

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// This interface will represent our car db
type PassportRepo interface {
	InsertInfo(ctx context.Context, info *PassportEntity) error
	GetInfo(ctx context.Context, name string) (*PassportEntity, error)
	RemoveInfo(ctx context.Context, name string) (*PassportEntity, error)
}

type passportRepoDeps struct {
	fx.In
}

func NewPassportRepo(deps passportRepoDeps) PassportRepo {
	return &passportMemoryRepo{
		deps:  deps,
		infos: make(map[string]*PassportEntity),
	}
}

type passportMemoryRepo struct {
	deps  passportRepoDeps
	infos map[string]*PassportEntity
}

func (c *passportMemoryRepo) InsertInfo(ctx context.Context, info *PassportEntity) error {
	if _, exists := c.infos[info.Name]; exists {
		return fmt.Errorf("car %s already exists", info.Name)
	}
	c.infos[info.Name] = info
	return nil
}

func (c *passportMemoryRepo) GetInfo(ctx context.Context, name string) (*PassportEntity, error) {
	if info, exists := c.infos[name]; exists {
		return info, nil
	}
	return nil, fmt.Errorf("unknown car ID %s", name)
}

func (c *passportMemoryRepo) RemoveInfo(ctx context.Context, name string) (*PassportEntity, error) {
	info, err := c.GetInfo(ctx, name)
	if err == nil {
		delete(c.infos, name)
		return info, nil
	}
	return nil, err
}
