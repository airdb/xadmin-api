package validations

import (
	"context"

	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
)

type PassportServiceValidations interface {
	Login(ctx context.Context, request *passportv1.LoginRequest) error
	Callback(ctx context.Context, request *passportv1.CallbackRequest) error
	Profile(ctx context.Context, request *passportv1.ProfileRequest) error
	Logout(ctx context.Context, request *passportv1.LogoutRequest) error
}

type passportServiceValidations struct {
}

func CreatePassportServiceValidations() PassportServiceValidations {
	return new(passportServiceValidations)
}

func (w *passportServiceValidations) Login(ctx context.Context, request *passportv1.LoginRequest) error {
	return nil
}

func (w *passportServiceValidations) Callback(ctx context.Context, request *passportv1.CallbackRequest) error {
	return nil
}

func (w *passportServiceValidations) Profile(ctx context.Context, request *passportv1.ProfileRequest) error {
	return nil
}

func (w *passportServiceValidations) Logout(ctx context.Context, request *passportv1.LogoutRequest) error {
	return nil
}
