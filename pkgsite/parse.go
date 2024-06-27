package pkgsite

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hhhapz/doc"
)

type ParseError struct {
	Sel     *goquery.Selection
	Message string
}

func (err ParseError) Error() string {
	return err.Message
}

type state struct {
	doc     *goquery.Document
	pkg     doc.Package
	current *doc.Type
	useCase bool
}

func newState(document *goquery.Document, useCase bool) (*state, error) {
	name := document.Find("h1.UnitHeader-titleHeading").Text()

	sel := document.Find("div.UnitDoc .Documentation-overview")
	overview := comments(sel.Children().NextUntil("details"))
	examples := examples(sel)
	url := document.Find("nav.go-Breadcrumb ol li").Last().Find("a").AttrOr("href", "/")[1:]

	subpkgs := subpackages(document)

	return &state{
		doc: document,
		pkg: doc.Package{
			URL:         url,
			Name:        name,
			Overview:    overview,
			Examples:    examples,
			ConstantMap: map[string]doc.Variable{},
			VariableMap: map[string]doc.Variable{},
			Functions:   map[string]doc.Function{},
			Types:       map[string]doc.Type{},
			Subpackages: subpkgs,
		},
		useCase: useCase,
	}, nil
}

func (s *state) newError(sel *goquery.Selection, msg string) error {
	return ParseError{sel, msg}
}

func (s *state) variables(sel *goquery.Selection, constants bool, m map[string]doc.Variable) error {
	sel.Filter(".Documentation-declaration").Each(func(i int, sel *goquery.Selection) {
		comment := comments(sel.NextUntil(".Documentation-declaration"))
		signature := sel.Find("pre").Text()
		v := doc.Variable{
			Signature: signature,
			Comment:   comment,
		}
		if constants {
			s.pkg.Constants = append(s.pkg.Constants, v)
		} else {
			s.pkg.Variables = append(s.pkg.Variables, v)
		}
		sel.Find("span[data-kind]").Each(func(i int, nameSel *goquery.Selection) {
			name := nameSel.AttrOr("id", "")
			v := doc.Variable{
				Name:      name,
				Signature: signature,
				Comment:   comment,
			}
			put(m, name, v, s.useCase)
		})
	})
	return nil
}

func (s *state) functions(sel *goquery.Selection) error {
	const base = ".Documentation-function"
	const header = "h4.Documentation-functionHeader a"

	sel.Each(func(i int, sel *goquery.Selection) {
		name := sel.Find(header).First().Text()
		decl := sel.Find("div.Documentation-declaration")
		comment := comments(decl.NextUntil("details"))
		f := doc.Function{
			Name:      name,
			Signature: strings.TrimSpace(decl.Text()),
			Comment:   comment,
		}
		put(s.pkg.Functions, name, f, s.useCase)
	})
	return nil
}

func (s *state) typ(sel *goquery.Selection) (doc.Type, error) {
	const header = "h4.Documentation-typeHeader a"

	name := sel.Find(header).First().Text()
	decl := sel.Find("div.Documentation-declaration").First()
	comment := comments(decl.NextUntil("details, .Documentation-typeFunc, .Documentation-typeMethod"))
	t := doc.Type{
		Name:          name,
		Signature:     strings.TrimSpace(decl.Text()),
		Comment:       comment,
		TypeFunctions: map[string]doc.Function{},
		Methods:       map[string]doc.Method{},
	}
	put(s.pkg.Types, name, t, s.useCase)
	return t, nil
}

func (s *state) typefuncs(sel *goquery.Selection, m map[string]doc.Function, dupe bool) error {
	const header = "h4.Documentation-typeFuncHeader a"

	name := sel.Find(header).First().Text()
	decl := sel.Find("div.Documentation-declaration").First()
	comment := comments(decl.NextUntil("details, .Documentation-typeFunc, .Documentation-typeMethod"))
	f := doc.Function{
		Name:      name,
		Signature: strings.TrimSpace(decl.Text()),
		Comment:   comment,
	}
	if dupe {
		put(s.pkg.Functions, name, f, s.useCase)
	}
	put(m, name, f, s.useCase)
	return nil
}

func (s *state) methods(sel *goquery.Selection, forType string, m map[string]doc.Method) error {
	const header = "h4.Documentation-typeMethodHeader a"

	name := sel.Find(header).First().Text()
	decl := sel.Find("div.Documentation-declaration").First()
	comment := comments(decl.NextUntil("details, .Documentation-typeFunc, .Documentation-typeMethod"))
	mtd := doc.Method{
		For: forType,
		Function: doc.Function{
			Name:      name,
			Signature: strings.TrimSpace(decl.Text()),
			Comment:   comment,
		},
	}
	put(m, name, mtd, s.useCase)
	return nil
}

const directoriesSelector = "h3#pkg-subdirectories"

func subpackages(doc *goquery.Document) []string {
	return nil
}

func comments(sel *goquery.Selection) doc.Comment {
	if sel.Length() == 0 {
		return nil
	}
	comments := make(doc.Comment, 0, len(sel.Nodes))
	sel.Each(func(i int, s *goquery.Selection) {
		n := s.Nodes[0]
		switch n.Data {
		case "p":
			f := strings.Fields(s.Text())
			comments = append(comments, doc.Paragraph(strings.Join(f, " ")))
		case "pre":
			comments = append(comments, doc.Pre(s.Text()))
		case "h4":
			if s.AttrOr("id", "") == "" {
				return
			}
			text := strings.TrimSpace(n.FirstChild.Data)
			comments = append(comments, doc.Heading(text))
		}
	})
	return comments
}

func examples(sel *goquery.Selection) []doc.Example {
	return nil
}

func put[V any](m map[string]V, name string, v V, useCase bool) {
	if !useCase {
		name = strings.ToLower(name)
	}
	m[name] = v
}
