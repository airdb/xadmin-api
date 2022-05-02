package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	apiv1 "github.com/airdb/xadmin-api/genproto/v1"
)

// FromProtoPassportInfoToModelPassportInfo converts workshop proto model to our data Entity
func FromProtoPassportInfoToModelPassportInfo(request *apiv1.LoginRequest) *data.PassportInfoEntity {
	if request == nil {
		return nil
	}
	return &data.PassportInfoEntity{
		Name:     request.GetName(),
		Password: request.GetPassword(),
	}
}

// FromModelPassportInfoToProtoPassportInfo converts our data Entity to workshop proto model
func FromModelPassportInfoToProtoPassportInfo(info *data.PassportInfoEntity) *apiv1.PassportInfo {
	if info == nil {
		return nil
	}
	return &apiv1.PassportInfo{
		Name: info.Name,
	}
}

// FromModelPassportInfoToSubWorkshopProtoPassportInfo converts our data Entity to workshop proto model
func FromModelPassportInfoToSubWorkshopMap(info *data.PassportInfoEntity) map[string]interface{} {
	if info == nil {
		return nil
	}

	return map[string]interface{}{
		"name": info.Name,
	}
}
