package querykit

import (
	"google.golang.org/protobuf/proto"
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
	GetSelect() *fieldmaskpb.FieldMask
	GetOmit() *fieldmaskpb.FieldMask
}

type Filter[TFilter any] interface {
	GetFilter() TFilter
}

type FieldMask struct {
	*fieldmaskpb.FieldMask
}

func NewField(f *fieldmaskpb.FieldMask) *FieldMask {
	return &FieldMask{
		f,
	}
}

func (f *FieldMask) GetSelect() *fieldmaskpb.FieldMask {
	fm, ok := proto.Clone(f.FieldMask).(*fieldmaskpb.FieldMask)
	if !ok {
		panic("can not clone mask message")
	}
	paths := []string{}
	for _, path := range fm.Paths {
		if len(path) == 0 {
			continue
		}
		if path[0] == '+' {
			paths = append(paths, path[1:])
		} else if path[0] == '-' {
			continue
		} else {
			paths = append(paths, path)
		}
	}
	return fm
}

func (f *FieldMask) GetOmit() *fieldmaskpb.FieldMask {
	fm, ok := proto.Clone(f.FieldMask).(*fieldmaskpb.FieldMask)
	if !ok {
		panic("can not clone mask message")
	}
	paths := []string{}
	for _, path := range fm.Paths {
		if len(path) == 0 {
			continue
		}
		if path[0] == '+' {
			continue
		} else if path[0] == '-' {
			paths = append(paths, path[1:])
		} else {
			continue
		}
	}
	return fm
}
