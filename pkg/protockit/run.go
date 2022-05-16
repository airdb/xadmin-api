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
	processors map[string]util.Processor
	actions    []string
}

func New(ctx context.Context, processors map[string]util.Processor) *Bootstrap {
	return &Bootstrap{
		ctx:        ctx,
		processors: processors,
		actions:    []string{},
	}
}

func (bs *Bootstrap) Run(gen *protogen.Plugin) error {
	params := util.ParseParameter(gen.Request.GetParameter())
	if _, ok := params["actions"]; ok {
		bs.actions = params["actions"]
	}
	if len(bs.actions) == 0 {
		for action, _ := range bs.processors {
			bs.actions = append(bs.actions, action)
		}
	}
	log.Println(bs.actions)

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

		for _, action := range bs.actions {
			if _, ok := bs.processors[action]; !ok {
				continue
			}
			if bs.ctx, err = bs.processors[action](bs.ctx, file); err != nil {
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
