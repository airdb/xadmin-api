package services

import (
	"context"

	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (w *uamImpl) CreateUser(ctx context.Context, request *uamv1.CreateUserRequest) (*uamv1.CreateUserResponse, error) {
	return w.deps.Controller.CreateUser(ctx, request)
}

func (w *uamImpl) GetUser(ctx context.Context, request *uamv1.GetUserRequest) (*uamv1.GetUserResponse, error) {
	return w.deps.Controller.GetUser(ctx, request)
}

func (w *uamImpl) ListUsers(ctx context.Context, request *uamv1.ListUsersRequest) (*uamv1.ListUsersResponse, error) {
	return w.deps.Controller.ListUsers(ctx, request)
}

func (w *uamImpl) UpdateUser(ctx context.Context, request *uamv1.UpdateUserRequest) (*uamv1.UpdateUserResponse, error) {
	return w.deps.Controller.UpdateUser(ctx, request)
}

func (w *uamImpl) DeleteUser(ctx context.Context, request *uamv1.DeleteUserRequest) (*emptypb.Empty, error) {
	return w.deps.Controller.DeleteUser(ctx, request)
}
