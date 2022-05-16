package querykit

import (
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Page handle pagination
type Page interface {
	GetPageOffset() int32
	GetPageSize() int32
}

// Search field
type Search interface {
	GetSearch() string
}

// Sort interface handle sorting like '+created_at','-created_at'
type Sort interface {
	GetSort() []string
}

// Fields fields to query or update
type Fields interface {
	Keep() *fieldmaskpb.FieldMask
	Omit() *fieldmaskpb.FieldMask
}

type Filter[TFilter any] interface {
	GetFilter() TFilter
}
