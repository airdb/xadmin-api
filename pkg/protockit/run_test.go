package protockit_test

import (
	"context"
	"log"
	"testing"

	. "github.com/airdb/xadmin-api/pkg/protockit"
	"github.com/airdb/xadmin-api/pkg/protockit/gencode"
	"github.com/airdb/xadmin-api/pkg/protockit/generror"
	"github.com/airdb/xadmin-api/pkg/protockit/genextends"
	"github.com/airdb/xadmin-api/pkg/protockit/tests"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/stretchr/testify/assert"
)

func TestRunTags(t *testing.T) {
	plugin, err := tests.NewPlugin("code.proto")
	assert.Nil(t, err)

	registedProcessort := map[string]util.Processor{
		"tags": generror.Process,
	}

	err = New(context.Background(), registedProcessort).Run(plugin)
	assert.Nil(t, err)
}

func TestRunError(t *testing.T) {
	plugin, err := tests.NewPlugin("library.proto")
	assert.Nil(t, err)

	registedProcessort := map[string]util.Processor{
		"error": generror.Process,
	}

	err = New(context.Background(), registedProcessort).Run(plugin)
	assert.Nil(t, err)
}

func TestRunExtends(t *testing.T) {
	plugin, err := tests.NewPlugin("library.proto")
	assert.Nil(t, err)

	registedProcessort := map[string]util.Processor{
		"extends": genextends.Process,
	}

	err = New(context.Background(), registedProcessort).Run(plugin)
	assert.Nil(t, err)
}

func TestRunCode(t *testing.T) {
	plugin, err := tests.NewPlugin("library.proto")
	assert.Nil(t, err)

	registedProcessort := map[string]util.Processor{
		"code": gencode.Process,
	}

	err = New(context.Background(), registedProcessort).Run(plugin)
	assert.Nil(t, err)

	response := plugin.Response()
	assert.Nil(t, response.Error)
	if response.Error != nil {
		log.Println(*response.Error)
	}

	assert.Greater(t, len(response.File), 0)

	if len(response.File) > 0 {
		log.Println(*response.File[0].Content)
	}
}
