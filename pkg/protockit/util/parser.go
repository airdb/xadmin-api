package util

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

const beginSymbol = `@Kit`

func ParseComment[T any](docs string) (T, error) {
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
