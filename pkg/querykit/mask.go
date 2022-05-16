package querykit

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	fmpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

type NamingMap map[string]string

func (nm NamingMap) Name(s string) string {
	if _, ok := nm[s]; ok {
		return nm[s]
	}
	return s
}

type KitFieldMask interface {
	KitFieldsByActions(action ...string) []string
	KitMaskMap(field string) string
}

type FieldMask struct {
	keep *fmpb.FieldMask
	omit *fmpb.FieldMask

	restrict []string
	kit      KitFieldMask
}

func NewField(f *fmpb.FieldMask, k KitFieldMask) *FieldMask {
	fm := &FieldMask{
		keep: f,
		omit: &fmpb.FieldMask{},
		kit:  k,
	}
	return fm
}

func (f *FieldMask) WithAction(action string) *FieldMask {
	if f.kit == nil {
		return f
	}
	f.restrict = f.kit.KitFieldsByActions(action)
	return f
}

func (f *FieldMask) Keep() *fmpb.FieldMask {
	fm, ok := proto.Clone(f.keep).(*fmpb.FieldMask)
	if !ok {
		panic("can not clone mask message")
	}

	if fm == nil {
		fm = &fmpb.FieldMask{
			Paths: append([]string{}, f.restrict...),
		}
	} else {
		paths := []string{}
		for _, path := range fm.Paths {
			if len(path) == 0 {
				continue
			}
			paths = append(paths, path)
		}
		fm.Paths = paths
	}
	fm.Paths = f.normalizer(lo.Intersect(f.restrict, fm.Paths)...)

	return fm
}

func (f *FieldMask) SetOmit(fm *fmpb.FieldMask) *FieldMask {
	f.omit = fm
	return f
}

func (f *FieldMask) AddOmit(s ...string) *FieldMask {
	f.omit.Paths = append(f.omit.Paths, s...)
	return f
}

func (f *FieldMask) Omit() *fmpb.FieldMask {
	fm, ok := proto.Clone(f.omit).(*fmpb.FieldMask)
	if !ok {
		panic("can not clone mask message")
	}
	paths := []string{}
	for _, path := range fm.Paths {
		if len(path) == 0 {
			continue
		}
		if len(path) == 0 {
			continue
		}
		paths = append(paths, path)
	}
	fm.Paths = f.normalizer(paths...)

	return fm
}

func (f *FieldMask) normalizer(fields ...string) []string {
	if f.kit == nil {
		return fields
	}

	res := []string{}
	for _, field := range fields {
		field = f.kit.KitMaskMap(field)
		if len(field) == 0 {
			continue
		}
		res = append(res, field)
	}

	return res
}
