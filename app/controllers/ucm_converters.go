package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	ucmv1 "github.com/airdb/xadmin-api/genproto/ucm/v1"
)

type ucmConvert struct{}

func newUcmConvert() *ucmConvert {
	return &ucmConvert{}
}

// User Convert Start

// FromProtoUserToModelUser converts proto model to our data Entity
func (c ucmConvert) FromProtoUserToModelUser(request *ucmv1.User) *data.UserEntity {
	if request == nil {
		return nil
	}
	return &data.UserEntity{
		Title: request.GetTitle(),
	}
}

// FromProtoUserToModelUser converts proto model to our data Entity
func (c ucmConvert) FromProtoCreateUserToModelUser(request *ucmv1.CreateUserRequest) *data.UserEntity {
	if request == nil {
		return nil
	}
	return &data.UserEntity{
		Title: request.User.GetTitle(),
	}
}

// FromModelUserToProtoUser converts our data Entity to proto model
func (c ucmConvert) FromModelUserToProtoUser(in *data.UserEntity) *ucmv1.User {
	if in == nil {
		return nil
	}

	return &ucmv1.User{
		Title: in.Title,
	}
}

// User Convert End
