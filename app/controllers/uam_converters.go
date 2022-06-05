package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
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
		Title: in.Title,
	}
}

// User Convert End
