package util

import (
	"context"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type utilContext int

const (
	dirCtx      = 1
	pluginCtx   = 2
	pluginFiles = 3
)

func NewContextDir(ctx context.Context, dir string) context.Context {
	return context.WithValue(ctx, dirCtx, dir)
}

func FromContextDir(ctx context.Context) string {
	if val, ok := ctx.Value(dirCtx).(string); ok {
		return val
	}
	return "."
}

func NewContextGen(ctx context.Context, gen *protogen.Plugin) context.Context {
	return context.WithValue(ctx, pluginCtx, gen)
}

func FromContextGen(ctx context.Context) *protogen.Plugin {
	if val, ok := ctx.Value(pluginCtx).(*protogen.Plugin); ok {
		return val
	}
	return nil
}

func NewContextFiles(ctx context.Context, files *protoregistry.Files) context.Context {
	return context.WithValue(ctx, pluginFiles, files)
}

func FromContextFiles(ctx context.Context) *protoregistry.Files {
	if val, ok := ctx.Value(pluginFiles).(*protoregistry.Files); ok {
		return val
	}
	return nil
}

func NewOrFromContextG(
	ctx context.Context, file *protogen.File,
) (context.Context, *protogen.GeneratedFile) {
	g := FromContextG(ctx, file)
	if g == nil {
		ctx = NewContextG(ctx, file)
	}
	return ctx, FromContextG(ctx, file)
}

func NewContextG(ctx context.Context, file *protogen.File) context.Context {
	gen := FromContextGen(ctx)
	if gen == nil {
		panic("plugin not in context")
	}

	g := FromContextG(ctx, file)
	if g != nil {
		return ctx
	}

	g = NewGeneratedFile(gen, file)
	return context.WithValue(ctx, file, g)
}

func FromContextG(ctx context.Context, file *protogen.File) *protogen.GeneratedFile {
	if val, ok := ctx.Value(file).(*protogen.GeneratedFile); ok {
		return val
	}
	return nil
}
