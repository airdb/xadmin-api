package controllers

import (
	"context"
	"errors"

	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"github.com/golang/protobuf/ptypes/empty"
)

func (c *uamController) ListUsers(ctx context.Context, request *uamv1.ListUsersRequest) (*uamv1.ListUsersResponse, error) {
	c.log.Debug(ctx, "list users accepted")

	total, filtered, err := c.deps.UserRepo.Count(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list users count error")
		return nil, errors.New("list users count error")
	}
	if total == 0 {
		return nil, errors.New("users is empty")
	}

	items, err := c.deps.UserRepo.List(ctx, request)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "list users error")
		return nil, errors.New("list users error")
	}

	return &uamv1.ListUsersResponse{
		TotalSize:    total,
		FilteredSize: filtered,
		Users: func() []*uamv1.User {
			res := make([]*uamv1.User, len(items))
			for i := 0; i < len(items); i++ {
				res[i] = c.conver.FromModelUserToProtoUser(items[i])
			}
			return res
		}(),
	}, nil
}

func (c *uamController) GetUser(ctx context.Context, request *uamv1.GetUserRequest) (*uamv1.GetUserResponse, error) {
	c.log.Debug(ctx, "get user accepted")

	item, err := c.deps.UserRepo.Get(ctx, request.GetId())
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get user error")
		return nil, errors.New("name not exist")
	}

	return &uamv1.GetUserResponse{
		User: c.conver.FromModelUserToProtoUser(item),
	}, err
}

func (c *uamController) CreateUser(ctx context.Context, request *uamv1.CreateUserRequest) (*uamv1.CreateUserResponse, error) {
	c.log.Debug(ctx, "create user accepted")

	item := c.conver.FromProtoCreateUserToModelUser(request)
	err := c.deps.UserRepo.Create(ctx, item)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "create user item failed")
		return nil, errors.New("create user item failed")
	}

	return &uamv1.CreateUserResponse{
		User: c.conver.FromModelUserToProtoUser(item),
	}, err
}

func (c *uamController) UpdateUser(ctx context.Context, request *uamv1.UpdateUserRequest) (*uamv1.UpdateUserResponse, error) {
	c.log.Debug(ctx, "update user accepted")
	data := c.conver.FromProtoUserToModelUser(request.GetUser())

	fm := querykit.NewField(request.GetUpdateMask(), request.GetUser()).WithAction("update")

	err := c.deps.UserRepo.Update(ctx, data.Id, data, fm)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "update user item failed")
		return nil, errors.New("update user item failed")
	}

	item, err := c.deps.UserRepo.Get(ctx, data.Id)
	if err != nil || item == nil {
		c.log.WithError(err).Debug(ctx, "update user item not exist")
		return nil, errors.New("update user item not exist")
	}

	return &uamv1.UpdateUserResponse{
		User: c.conver.FromModelUserToProtoUser(item),
	}, err
}

func (c *uamController) DeleteUser(ctx context.Context, request *uamv1.DeleteUserRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete user accepted")

	err := c.deps.UserRepo.Delete(ctx, request.GetId())
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete user item failed")
		return nil, errors.New("delete user item failed")
	}

	return &empty.Empty{}, nil
}
