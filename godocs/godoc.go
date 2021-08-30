package godocs

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hhhapz/doc"
)

const (
	// base is the base path to godocs.io
	base = "https://godocs.io/"
	// selectors matches methods, types, and functions on the documentation
	// page.
	selectors = `[data-kind="function"], [data-kind="type"], [data-kind="method"]:not([class*="decl"])`
)

// godocParser implements doc.Parser.
type godocParser struct{}

// Parser is an implementation of godoc.Parser that retrieves documentation
// from https://godocs.io.
var Parser doc.Parser = godocParser{}

// URL returns a url to the path to see the documentation for the provided
// module on https://godocs.io/.
func (godocParser) URL(module string) string {
	return base + module
}

func (p godocParser) Parse(document *goquery.Document) (doc.Package, error) {
	// special case not found case for godocs
	if document.Find("head title").Text() == "Not Found - godocs.io" {
		return doc.Package{}, doc.InvalidStatusError(404)
	}

	s := newState(document)

	var err error
	document.Find(selectors).EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		kind := sel.AttrOr("data-kind", "")
		switch kind {
		case "function":
			err = s.function(sel)
		case "type":
			err = s.typ(sel)
		case "method":
			err = s.method(sel)
		default:
			// this should never happen
		}
		// true when err is nil
		return err == nil
	})
	s.pkg.Types[s.current.Name] = *s.current

	if err != nil {
		return doc.Package{}, err
	}
	return s.pkg, nil
}
