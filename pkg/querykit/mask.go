package querykit

import (
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

type Nameing interface {
	func(string) string | map[string]string
}

type FieldMask struct {
	keep *fmpb.FieldMask
	omit *fmpb.FieldMask

	naming func(string) string
}

func NewField[T Nameing](f *fmpb.FieldMask, n T) *FieldMask {
	fm := &FieldMask{
		keep: f,
		omit: &fmpb.FieldMask{},
	}
	func(n any) {
		switch nt := n.(type) {
		case func(string) string:
			fm.naming = nt
		case map[string]string:
			fm.naming = NamingMap(nt).Name
		}
	}(n)

	return fm
}

func (f *FieldMask) Keep() *fmpb.FieldMask {
	fm, ok := proto.Clone(f.keep).(*fmpb.FieldMask)
	if !ok {
		panic("can not clone mask message")
	}
	paths := []string{}
	for _, path := range fm.Paths {
		if len(path) == 0 {
			continue
		}
		path = f.normalizer(path)
		if len(path) == 0 {
			continue
		}
		paths = append(paths, path)
	}
	fm.Paths = paths
	return fm
}

func (f *FieldMask) SetOmit(fm *fmpb.FieldMask) {
	f.omit = fm
}

func (f *FieldMask) AddOmit(s ...string) {
	f.omit.Paths = append(f.omit.Paths, s...)
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
		path = f.normalizer(path)
		if len(path) == 0 {
			continue
		}
		paths = append(paths, path)
	}
	fm.Paths = paths
	return fm
}

func (f *FieldMask) normalizer(s string) string {
	if f.naming == nil {
		return s
	}

	return f.naming(s)
}
