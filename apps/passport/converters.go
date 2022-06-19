package passport

import (
	"github.com/airdb/xadmin-api/apps/data"
	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
	"github.com/casdoor/casdoor-go-sdk/auth"
)

type passportConvert struct{}

// FromProtoInfoToModelInfo converts proto model to our data Entity
func (c passportConvert) FromProtoInfoToModelInfo(request *passportv1.LoginRequest) *data.PassportEntity {
	if request == nil {
		return nil
	}
	return &data.PassportEntity{
		Name:     request.GetName(),
		Password: request.GetPassword(),
	}
}

// FromModelInfoToProtoInfo converts our data Entity to proto model
func (c passportConvert) FromModelInfoToProtoInfo(info *data.PassportEntity) *passportv1.Info {
	if info == nil {
		return nil
	}
	return &passportv1.Info{
		Name: info.Name,
	}
}

// FromModelInfoToProtoInfo converts our data Entity to proto model
func (c passportConvert) FromClaimsUserToProtoInfo(user *auth.User) *passportv1.Info {
	if user == nil {
		return nil
	}

	return &passportv1.Info{
		Id:          user.Id,
		Name:        user.Name,
		Type:        user.Type,
		DisplayName: user.DisplayName,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		IsAdmin:     user.IsAdmin,
		IsForbidden: user.IsForbidden,
		IsDeleted:   user.IsDeleted,
	}
}
