package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/airdb/xadmin-api/pkg/protockit"
	"github.com/airdb/xadmin-api/pkg/protockit/gentags"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/airdb/xadmin-api/pkg/protockit/version"
	"google.golang.org/protobuf/compiler/protogen"
)

var flagVersion bool
var registedProcessort map[string]util.Processor

func init() {
	// log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.BoolVar(&flagVersion, "version", false, "print plugin version")

	registedProcessort = map[string]util.Processor{
		"tags": gentags.Process,
	}
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Printf("%v %v\n", filepath.Base(os.Args[0]), version.PluginVersion)
		os.Exit(0)
	}

	bs := protockit.New(context.Background(), registedProcessort)
	protogen.Options{}.Run(bs.Run)
}
