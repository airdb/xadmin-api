package tests

import (
	"github.com/pseudomuto/protokit/utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewPlugin(file string) (*protogen.Plugin, error) {
	set, err := utils.LoadDescriptorSet("./fixtures", "fileset.pb")
	if err != nil {
		return nil, err
	}

	req := utils.CreateGenRequest(set, file)
	req.Parameter = proto.String("paths=source_relative")

	options := protogen.Options{}
	plugin, err := options.New(req)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}
