package validations

import (
	"context"

	apiv1 "github.com/airdb/xadmin-api/genproto/v1"
)

type PassportServiceValidations interface {
	Login(ctx context.Context, request *apiv1.LoginRequest) error
	LoginCallback(ctx context.Context, request *apiv1.LoginCallbackRequest) error
	Profile(ctx context.Context, request *apiv1.ProfileRequest) error
	Logout(ctx context.Context, request *apiv1.LogoutRequest) error
}

type passportServiceValidations struct {
}

func CreatePassportServiceValidations() PassportServiceValidations {
	return new(passportServiceValidations)
}

func (w *passportServiceValidations) Login(ctx context.Context, request *apiv1.LoginRequest) error {
	return nil
}

func (w *passportServiceValidations) LoginCallback(ctx context.Context, request *apiv1.LoginCallbackRequest) error {
	return nil
}

func (w *passportServiceValidations) Profile(ctx context.Context, request *apiv1.ProfileRequest) error {
	return nil
}

func (w *passportServiceValidations) Logout(ctx context.Context, request *apiv1.LogoutRequest) error {
	return nil
}
