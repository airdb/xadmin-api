package data

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// This interface will represent our car db
type PassportInfoDB interface {
	InsertPassportInfo(ctx context.Context, info *PassportInfoEntity) error
	GetPassportInfo(ctx context.Context, name string) (*PassportInfoEntity, error)
	RemovePassportInfo(ctx context.Context, name string) (*PassportInfoEntity, error)
}

type passportInfoDBDeps struct {
	fx.In
}

func CreatePassportInfoDB(deps passportInfoDBDeps) PassportInfoDB {
	return &passportInfoDB{
		deps:  deps,
		infos: make(map[string]*PassportInfoEntity),
	}
}

type passportInfoDB struct {
	deps  passportInfoDBDeps
	infos map[string]*PassportInfoEntity
}

func (c *passportInfoDB) InsertPassportInfo(ctx context.Context, info *PassportInfoEntity) error {
	if _, exists := c.infos[info.Name]; exists {
		return fmt.Errorf("car %s already exists", info.Name)
	}
	c.infos[info.Name] = info
	return nil
}

func (c *passportInfoDB) GetPassportInfo(ctx context.Context, name string) (*PassportInfoEntity, error) {
	if info, exists := c.infos[name]; exists {
		return info, nil
	}
	return nil, fmt.Errorf("unknown car ID %s", name)
}

func (c *passportInfoDB) RemovePassportInfo(ctx context.Context, name string) (*PassportInfoEntity, error) {
	info, err := c.GetPassportInfo(ctx, name)
	if err == nil {
		delete(c.infos, name)
		return info, nil
	}
	return nil, err
}
