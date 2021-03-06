package gentags

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/gobeam/stringy"
	"google.golang.org/protobuf/compiler/protogen"
)

// tagExp will find a single tag expressions
var tagExp = regexp.MustCompile(`[_a-z][_\w]*:".+?"`)

func Process(ctx context.Context, file *protogen.File) (context.Context, error) {
	dir := util.FromContextDir(ctx)

	fileinfo, err := os.Stat(dir)
	if err != nil {
		log.Printf("could not open file or directory: %s, %v\n", dir, err)
		return ctx, err
	}

	var parseError error

	// enable parsing of single files or directories
	switch fileinfo.IsDir() {

	case false:
		// handle single files
		parseError = newWalker(ctx, file).applyStructTags(dir)

	case true:
		// handle all files in the given directory
		parseError = newWalker(ctx, file).walk(dir)
	}

	if parseError != nil {
		log.Printf("error while parsing files: %s, %v\n", dir, err)
	}

	return ctx, nil
}

type walker struct {
	ctx    context.Context
	file   *protogen.File
	suffix string
}

func newWalker(ctx context.Context, file *protogen.File) *walker {
	suffix := file.GeneratedFilenamePrefix
	if len(suffix) == 0 {
		suffix = ".pb.go"
	}

	if strings.HasSuffix(suffix, ".pb.go") == false {
		suffix += ".pb.go"
	}

	return &walker{
		ctx:    ctx,
		file:   file,
		suffix: suffix,
	}
}

func (w *walker) walk(dir string) error {
	return filepath.Walk(dir, w.walker)
}

// walker will be executed recusrively on every file in the target directory
func (w *walker) walker(path string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	path = filepath.ToSlash(path)

	if strings.HasPrefix(path, ".git") {
		return nil
	}

	// ignore all files that have not been generated by pb
	if strings.HasSuffix(path, w.suffix) == false {
		return nil
	}

	err = w.applyStructTags(path)
	if err != nil {
		log.Printf("%s: %v\n", path, err)
	}

	return err
}

// applyStructTags will add specifically defined tags to the go protobuf structs
func (w *walker) applyStructTags(path string) error {

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("could not parse go file: %v", err)
	}

	var visitor visitor
	ast.Walk(visitor, astFile)

	gen := util.FromContextGen(w.ctx)
	if gen == nil {
		return fmt.Errorf("could not get gen in context")
	}

	g := gen.NewGeneratedFile(w.file.GeneratedFilenamePrefix+".pb.go", w.file.GoImportPath)

	// file, err := os.Create(path)
	// if err != nil {
	// 	return fmt.Errorf("could not open output file: %v", err)
	// }
	// defer file.Close() // nolint: errcheck

	err = format.Node(g, fileSet, astFile)
	if err != nil {
		return fmt.Errorf("could not write output file: %v", err)
	}

	return nil
}

// define a visitor struct to handle every ast node
type visitor struct{}

// Visit will find all fields in the code and add tags that are placed
// in the comment above the field as struct tags to the field
func (v visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch t := node.(type) {

	// we are only interested in field types
	case *ast.Field:

		// comments after the field are stripped away by the grpc protoc-generator
		// therefore we need to check the lines above the field which are
		// contained in the comment that is part of the field
		var commentTags []string
		docs := bytes.NewBuffer(nil)

		// iterate through the sub-nodes of the field
		ast.Inspect(node, func(child ast.Node) bool {

			// handle each comment line separately
			// note: commentGroup could be used to handle all comments at once
			switch comment := child.(type) {
			case *ast.Comment:

				// we are only interested in comments which matches our specific syntax
				_, err := fmt.Fprintln(docs, comment.Text)
				if err != nil {
					return false
				}
			}

			return true

		})

		// parse comment
		opt, err := util.KitParser[*annov1.FieldDescriptor](docs.String())
		if err != nil {
			break
		}

		if opt == nil {
			opt = &annov1.FieldDescriptor{
				Tags: map[string]string{},
			}
		}
		if opt.Tags == nil {
			opt.Tags = map[string]string{}
		}

		// add default yaml tag
		if _, ok := opt.Tags["yaml"]; !ok && len(t.Names) > 0 {
			opt.Tags["yaml"] = stringy.New(t.Names[0].Name).LcFirst()
		}

		// add option tags to commentTags
		for k, v := range opt.Tags {
			if len(v) == 0 {
				continue
			}
			commentTags = append(commentTags, fmt.Sprintf(`%s:"%s"`, k, v))
		}

		// continue with the next node if there are no comment tags
		if len(commentTags) == 0 {
			break
		}

		// tags are stored as child nodes of type ast.BasicLit
		ast.Inspect(node, func(child ast.Node) bool {

			switch basicLit := child.(type) {
			case *ast.BasicLit:

				// combine the current tag list with our new comment tags
				basicLit.Value = combineTags(basicLit.Value, commentTags)
			}

			return true
		})

	}

	return v
}

// combineTags will combine the given tags with the existing tags
// the order will be kept as is with existing tags first and new tags after.
// new tags will replace existing tags. if a tag is repeated, only the last
// version will be kept
func combineTags(tagList string, commentTags []string) string {

	// remove ticks from the current tag list
	tagList = strings.Trim(tagList, "`")

	// commentTags may contain multiple struct tags i.e. db:"smthing" json:"test"
	// therefore we combine all comment tags and split them into separate
	// tags again
	commentTagList := strings.Join(commentTags, " ")

	// use a regular expression to split all comments into tags
	tags := tagExp.FindAllString(commentTagList, -1)

	for _, tag := range tags {

		// get the tag name from the complete tag
		tagName := strings.Split(tag, ":")[0]

		// define the tag that should be matched
		tagMatch := fmt.Sprintf(`%s:".+?"`, tagName)

		// create a custom tag matcher based on the tag name
		matcher := regexp.MustCompile(tagMatch)

		// replace any matching existing tags in the taglist
		if matcher.MatchString(tagList) {
			tagList = matcher.ReplaceAllString(tagList, tag)
			continue
		}

		// append the tag to the list of tags
		tagList = fmt.Sprintf("%s %s", tagList, tag)

	}

	return fmt.Sprintf("`%s`", tagList)
}
