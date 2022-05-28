package util

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const beginSymbol = `@Kit`

type GetGencodeInterface interface {
	GetGencode() *annov1.GencodeOption
}

func KitParser[T any](docs string) (T, error) {
	var stared bool
	buf := bytes.NewBufferString(docs)
	data := bytes.NewBuffer(nil)
	for {
		line, err := buf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if !stared && strings.TrimRight(line[3:], "\n") == beginSymbol {
			stared = true
			continue
		}
		if stared && len(strings.TrimRight(line[2:], "\n")) == 0 {
			stared = false
			break
		}
		if stared == true {
			data.WriteString(line[4:])
		}
	}

	var dest T
	err := yaml.Unmarshal(data.Bytes(), &dest)

	return dest, err
}

func KitGencodeValid(message GetGencodeInterface) bool {
	return message.GetGencode() != nil && message.GetGencode().Layers != nil
}

func KitGencodeLayerEmpty(message GetGencodeInterface) bool {
	return KitGencodeValid(message) && len(message.GetGencode().Layers) == 0
}

func KitGencodeLayerValid(message GetGencodeInterface, layers ...string) bool {
	return KitGencodeValid(message) &&
		len(lo.Intersect(message.GetGencode().Layers, layers)) > 0
}

func ParseParameter(s string) url.Values {
	res := url.Values{}
	for _, item := range strings.Split(s, ",") {
		kvp := strings.SplitN(item, "=", 2)
		if len(kvp) != 2 {
			continue
		}
		res.Add(kvp[0], kvp[1])
	}

	return res
}
