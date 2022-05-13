package protockit

import (
	"context"
	"log"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Bootstrap struct {
	ctx        context.Context
	processors []util.Processor
}

func New(ctx context.Context, processors ...util.Processor) *Bootstrap {
	return &Bootstrap{
		ctx:        ctx,
		processors: processors,
	}
}

func (bs *Bootstrap) Run(gen *protogen.Plugin) error {
	var files protoregistry.Files
	for _, file := range gen.Files {
		if err := files.RegisterFile(file.Desc); err != nil {
			return err
		}
	}

	bs.ctx = util.NewContextGen(bs.ctx, gen)
	bs.ctx = util.NewContextFiles(bs.ctx, &files)

	var err error
	for _, file := range gen.Files {
		file := file
		if !file.Generate {
			continue
		}

		// gentags.Process(".", file.GeneratedFilenamePrefix)

		// if err := gencode.Process(g, file, &files); err != nil {
		// 	return err
		// }
		// if err := genextends.Process(g, file, &files); err != nil {
		// 	return err
		// }
		for _, processor := range bs.processors {
			if bs.ctx, err = processor(bs.ctx, file); err != nil {
				return err
			}

		}
	}
	return nil
}

func printContent(g *protogen.GeneratedFile) {
	// for debug
	if content, err := g.Content(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s\n", content)
	}
}
