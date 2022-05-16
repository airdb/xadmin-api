package protockit_test

import (
	"context"
	"testing"

	. "github.com/airdb/xadmin-api/pkg/protockit"
	"github.com/airdb/xadmin-api/pkg/protockit/generror"
	"github.com/airdb/xadmin-api/pkg/protockit/genextends"
	"github.com/airdb/xadmin-api/pkg/protockit/tests"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
)

func TestRunCode(t *testing.T) {
	plugin, err := tests.NewPlugin("code.proto")
	require.NoError(t, err)

	err = New(context.Background(), map[string]util.Processor{
		"test": func(ctx context.Context, file *protogen.File) (context.Context, error) {
			return ctx, nil
		},
	}).Run(plugin)
	require.NoError(t, err)
}

func TestRunLibrary(t *testing.T) {
	plugin, err := tests.NewPlugin("library.proto")
	require.NoError(t, err)

	registedProcessort := map[string]util.Processor{
		"error":   generror.Process,
		"extends": genextends.Process,
	}

	err = New(context.Background(), registedProcessort).Run(plugin)
	require.NoError(t, err)
}
