package util

import (
	"context"

	"google.golang.org/protobuf/compiler/protogen"
)

type Processor func(ctx context.Context, file *protogen.File) (context.Context, error)
