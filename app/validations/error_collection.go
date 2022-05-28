package validations

import (
	"context"

	errorCollectionv1 "github.com/airdb/xadmin-api/genproto/error_collection/v1"
)

type ErrorCollectionServiceValidations interface {
	Collect(context.Context, *errorCollectionv1.CreateErrorCollectionRequest) error
}

type errorCollectionServiceValidations struct {
}

func CreateErrorCollectionServiceValidations() ErrorCollectionServiceValidations {
	return new(errorCollectionServiceValidations)
}

func (w *errorCollectionServiceValidations) Collect(context.Context, *errorCollectionv1.CreateErrorCollectionRequest) error {
	return nil
}
