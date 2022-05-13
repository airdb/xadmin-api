package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/airdb/xadmin-api/pkg/protockit"
	"github.com/airdb/xadmin-api/pkg/protockit/gencode"
	"github.com/airdb/xadmin-api/pkg/protockit/genextends"
	"github.com/airdb/xadmin-api/pkg/protockit/version"
	"google.golang.org/protobuf/compiler/protogen"
)

var flagVersion bool
var flagActions string

func init() {
	// log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.BoolVar(&flagVersion, "version", false, "print plugin version")
	flag.StringVar(&flagActions, "actions", "tags,code,extends", "the selected actions, options: tags, code, extends")
}

func main() {
	flag.Parse()
	if flagVersion {
		log.Printf("%v %v\n", filepath.Base(os.Args[0]), version.PluginVersion)
		os.Exit(0)
	}

	bs := protockit.New(context.Background(), gencode.Process, genextends.Process)
	protogen.Options{}.Run(bs.Run)
}
