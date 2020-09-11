package generator

//go:generate statik -src templates

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/rakyll/statik/fs"

	// embed templates
	_ "github.com/bots-house/go-text-gen/generator/statik"

	"github.com/bots-house/go-text-gen/loader"
	"github.com/iancoleman/strcase"
)

type ctx struct {
	Bundle           *loader.Bundle
	PkgName          string
	TypeLocaleName   string
	TypeMessagesName string
}

func loadEmbeddedTemplate() ([]byte, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	r, err := statikFS.Open("/template.go.tmpl")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

// Generate Go file from bundle.
// If template is nil, use default.
func Generate(
	bundle *loader.Bundle,
	tmpl []byte,
	pkgName string,
) ([]byte, error) {
	if tmpl == nil {
		var err error
		tmpl, err = loadEmbeddedTemplate()
		if err != nil {
			return nil, fmt.Errorf("load embedded template: %w", err)
		}
	}

	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"upper":      strings.ToUpper,
		"camel":      strcase.ToCamel,
		"lowerCamel": strcase.ToLowerCamel,
	}).Parse(string(tmpl))

	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	buf := &bytes.Buffer{}

	if err := t.Execute(buf, ctx{
		Bundle:           bundle,
		PkgName:          pkgName,
		TypeLocaleName:   "Locale",
		TypeMessagesName: "Messages",
	}); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return format.Source(buf.Bytes())
}
