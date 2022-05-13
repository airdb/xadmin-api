package genextends

import (
	"context"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func Process(ctx context.Context, file *protogen.File) (context.Context, error) {
	var rangeErr error
	ctx, g := util.NewOrFromContextG(ctx, file)
	for _, message := range file.Messages {
		RangeFieldDescriptorsInMessage(
			message,
			func(option *annov1.MessageDescriptor, options map[string]*annov1.FieldDescriptor) bool {
				err := NewMessageGenerator(message, option, options).Run(g)
				if err != nil {
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

// RangeFieldDescriptorsInMessage iterates over all fileds in a message while fn returns true.
// The iteration order is undefined.
func RangeFieldDescriptorsInMessage(
	message *protogen.Message,
	fn func(option *annov1.MessageDescriptor, options map[string]*annov1.FieldDescriptor) bool,
) {
	if len(message.Fields) == 0 {
		return
	}

	option, err := util.ParseComment[*annov1.MessageDescriptor](
		message.Comments.Leading.String())
	if err != nil || option == nil {
		option = &annov1.MessageDescriptor{}
	}
	if option.DefaultFieldActions == nil {
		option.DefaultFieldActions = []string{"*"}
	}

	opts := map[string]*annov1.FieldDescriptor{}
	for _, field := range message.Fields {
		opt, err := util.ParseComment[*annov1.FieldDescriptor](
			field.Comments.Leading.String())
		if err != nil || opt == nil {
			opt = &annov1.FieldDescriptor{}
		}
		if opt.Actions == nil || len(opt.Actions) == 0 {
			opt.Actions = append([]string{}, option.DefaultFieldActions...)
		}
		opts[field.GoIdent.GoName] = opt
	}

	if !fn(option, opts) {
		return
	}
}
