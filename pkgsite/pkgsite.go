package pkgsite

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hhhapz/doc"
)

const (
	// base is the base path to godocs.io
	base = "https://pkg.go.dev/"
)

// pkgsiteParser implements doc.Parser.
type pkgsiteParser struct{}

// Parser is an implementation of godoc.Parser that retrieves documentation
// from https://godocs.io.
var Parser doc.Parser = pkgsiteParser{}

// URL returns a url to the path to see the documentation for the provided
// module on https://godocs.io/.
func (pkgsiteParser) URL(module string) string {
	return base + module
}

func (p pkgsiteParser) Parse(document *goquery.Document, useCase bool) (doc.Package, error) {
	// special case not found case for godocs
	if document.Find("h3.Error-message").Text() == "404 Not Found" {
		return doc.Package{}, doc.InvalidStatusError(404)
	}

	s, err := newState(document, useCase)
	if err != nil {
		return doc.Package{}, err
	}

	consts := document.Find("section.Documentation-constants")
	s.variables(consts.Children(), true, s.pkg.ConstantMap)

	vars := document.Find("section.Documentation-variables")
	s.variables(vars.Children(), false, s.pkg.VariableMap)

	funcs := document.Find(".Documentation-function")
	s.functions(funcs)

	types := document.Find(".Documentation-type")
	types.Each(func(i int, sel *goquery.Selection) {
		t, _ := s.typ(sel)
		sel.Find(".Documentation-typeFunc").Each(func(i int, sel *goquery.Selection) {
			s.typefuncs(sel, t.TypeFunctions)
		})
		sel.Find(".Documentation-typeMethod").Each(func(i int, sel *goquery.Selection) {
			s.methods(sel, t.Name, t.Methods)
		})
	})

	return s.pkg, nil
}
