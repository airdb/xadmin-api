package generror

import (
	"context"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func Process(ctx context.Context, file *protogen.File) (context.Context, error) {
	var rangeErr error
	ctx, g := util.NewOrFromContextG(ctx, file)
	for _, enum := range file.Enums {
		RangeEnumCodeDescriptorsInFile(
			enum,
			func(
				enum *protogen.Enum,
				option *annov1.ErrorEnumCodeDescriptor,
				opts map[string]*annov1.ErrorEnumValueDescriptor) bool {
				g.Unskip()
				if err := NewGenerator(enum, option, opts).Run(g); err != nil {
					rangeErr = err
					return false
				}
				return true
			},
		)
	}
	if rangeErr != nil {
		return ctx, rangeErr
	}
	return ctx, nil
}

// RangeEnumCodeDescriptorsInFile iterates over all resource descriptors in a file while fn returns true.
// The iteration order is undefined.
func RangeEnumCodeDescriptorsInFile(
	enum *protogen.Enum,
	fn func(
		enum *protogen.Enum,
		option *annov1.ErrorEnumCodeDescriptor,
		opts map[string]*annov1.ErrorEnumValueDescriptor) bool,
) {
	option, err := util.KitParser[*annov1.ErrorEnumCodeDescriptor](
		enum.Comments.Leading.String())
	if err != nil {
		option = &annov1.ErrorEnumCodeDescriptor{}
	}
	if option.Category != "error" {
		return
	}
	if option.DefaultHttpCode == 0 {
		option.DefaultHttpCode = 500
	}
	if len(option.DefaultDescribe) == 0 {
		option.DefaultDescribe = "error"
	}

	opts := map[string]*annov1.ErrorEnumValueDescriptor{}
	for _, value := range enum.Values {
		opt, err := util.KitParser[*annov1.ErrorEnumValueDescriptor](
			value.Comments.Leading.String())
		if err != nil {
			panic(err)
		}
		if opt.HttpCode == 0 {
			opt.HttpCode = option.DefaultHttpCode
		}
		if len(opt.Describe) == 0 {
			opt.Describe = option.DefaultDescribe
		}
		opts[value.GoIdent.GoName] = opt
	}

	if !fn(enum, option, opts) {
		return
	}
}
