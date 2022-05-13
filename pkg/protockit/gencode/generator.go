package gencode

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/compiler/protogen"
	// "github.com/stoewer/go-strcase"
	// "go.einride.tech/aip/reflect/aipreflect"
	// "go.einride.tech/aip/resourcename"
	// "google.golang.org/protobuf/compiler/protogen"
	// "google.golang.org/protobuf/reflect/protoreflect"
	// "google.golang.org/protobuf/reflect/protoregistry"
)

const (
	errorsPackage = protogen.GoImportPath("github.com/airdb/xadmin-api/pkg/errorskit")
	fmtPackage    = protogen.GoImportPath("fmt")
)

type generator struct {
	enum   *protogen.Enum
	option *annov1.ErrorEnumCodeDescriptor
	opts   map[string]*annov1.ErrorEnumValueDescriptor
}

func NewGenerator(
	enum *protogen.Enum,
	option *annov1.ErrorEnumCodeDescriptor,
	opts map[string]*annov1.ErrorEnumValueDescriptor,
) generator {
	return generator{
		enum:   enum,
		option: option,
		opts:   opts,
	}
}

func (r generator) Run(g *protogen.GeneratedFile) error {
	if r.enum == nil || r.option == nil || len(r.option.Type) == 0 {
		return nil
	}

	if len(r.enum.Values) > 0 {
		if err := r.registCode(g); err != nil {
			return err
		}
	}

	g.QualifiedGoIdent(errorsPackage.Ident(""))
	g.QualifiedGoIdent(fmtPackage.Ident(""))

	return nil
}

func (r generator) registCode(g *protogen.GeneratedFile) error {
	enumName := r.enum.GoIdent.GoName

	enumCodeType := fmt.Sprintf("ERRORCODE_TYPE_%s", enumName)
	g.P(fmt.Sprintf("const %s = \"%s\"", enumCodeType, r.option.Type))
	g.P()
	enumCodeOffset := fmt.Sprintf("ERRORCODE_OFFSET_%s", enumName)
	g.P(fmt.Sprintf("const %s = %d", enumCodeOffset, r.option.Offset))
	g.P()

	g.P(fmt.Sprintf("func (x %s) EnumType() string {", enumName))
	g.P(fmt.Sprintf(`return %s`, enumCodeType))
	g.P("}")
	g.P()
	g.P(fmt.Sprintf("func (x %s) EnumOffset() int {", enumName))
	g.P(fmt.Sprintf(`return %s`, enumCodeOffset))
	g.P("}")
	g.P()
	g.P(fmt.Sprintf("func (x %s) OffsetCode() int {", enumName))
	g.P("return x.EnumOffset() + int(x)")
	g.P("}")

	g.P()
	g.P(fmt.Sprintf("func Regist%s(fn func(code int, httpCode int, msg string)) {", enumName))
	var vis valueInfos
	for i := 0; i < len(r.enum.Values); i++ {
		value := r.enum.Values[i]
		originalName := value.GoIdent.GoName
		opt, ok := r.opts[value.GoIdent.GoName]
		if !ok {
			continue
		}
		g.P(fmt.Sprintf(`fn(%s.OffsetCode(), %d, "%s")`, originalName, opt.HttpCode, opt.Describe))
		vi := &valueInfo{
			Name:       string(r.enum.Desc.Name()),
			Value:      string(value.Desc.Name()),
			CamelValue: case2Camel(string(value.Desc.Name())),
			HTTPCode:   opt.HttpCode,
		}
		vis = append(vis, vi)
	}
	g.P("}")

	g.P()
	g.P(vis.execute())
	return nil
}

// ParseComment return code and msg
func (r generator) ParseComment(originalName, comment string) (int, string) {
	reg := regexp.MustCompile(`\w*\s*-\s*(\d{3})\s*:\s*(\w.*)\s*\.*`)
	if !reg.MatchString(comment) {
		log.Printf("constant '%s' have wrong comment(%s) format, register with 500 as default", originalName, comment)

		return 500, "Internal server error"
	}

	groups := reg.FindStringSubmatch(comment)
	if len(groups) != 3 {
		return 500, "Internal server error"
	}

	code, err := strconv.Atoi(groups[1])
	if err != nil {
		panic(fmt.Errorf("error number: %s", groups[1]))
	}

	return code, groups[2]
}

var caseTitle = cases.Title(language.Chinese)

func case2Camel(name string) string {
	if !strings.Contains(name, "_") {
		upperName := strings.ToUpper(name)
		if upperName == name {
			name = strings.ToLower(name)
		}
		// return cases.Title(language.Chinese)
		return caseTitle.String(name)
	}
	name = strings.ToLower(name)
	name = strings.Replace(name, "_", " ", -1)
	name = caseTitle.String(name)
	return strings.Replace(name, " ", "", -1)
}
