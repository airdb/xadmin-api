package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type uamConvert struct{}

func newUamConvert() *uamConvert {
	return &uamConvert{}
}

// User Convert Start

// FromProtoUserToModelUser converts proto model to our data Entity
func (c uamConvert) FromProtoUserToModelUser(request *uamv1.User) *data.UserEntity {
	if request == nil {
		return nil
	}
	return &data.UserEntity{
		Title: request.GetTitle(),
	}
}

// FromProtoUserToModelUser converts proto model to our data Entity
func (c uamConvert) FromProtoCreateUserToModelUser(request *uamv1.CreateUserRequest) *data.UserEntity {
	if request == nil {
		return nil
	}
	return &data.UserEntity{
		Title: request.User.GetTitle(),
	}
}

// FromModelUserToProtoUser converts our data Entity to proto model
func (c uamConvert) FromModelUserToProtoUser(in *data.UserEntity) *uamv1.User {
	if in == nil {
		return nil
	}

	return &uamv1.User{
		Id:        in.Id,
		Type:      in.Type,
		CreatedAt: timestamppb.New(*in.CreatedAt.Time),
		UpdatedAt: func() *timestamppb.Timestamp {
			if in.UpdatedAt.Time == nil {
				return nil
			}
			return timestamppb.New(*in.UpdatedAt.Time)
		}(),
		Username:      in.Username,
		DisplayName:   in.DisplayName,
		Avatar:        in.Avatar,
		Email:         in.Email,
		Phone:         in.Phone,
		Title:         in.Title,
		IsOnline:      in.IsOnline,
		IsAdmin:       in.IsAdmin,
		IsGlobalAdmin: in.IsGlobalAdmin,
		IsForbidden:   in.IsForbidden,
		IsDeleted:     in.IsDeleted,
	}
}

// User Convert End
