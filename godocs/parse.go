package godocs

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hhhapz/doc"
	"golang.org/x/net/html"
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

func newState(document *goquery.Document, useCase bool) *state {
	name := document.Find("#pkg-overview").Text()
	name = strings.TrimPrefix(name, "package ")

	sel := document.Find("#pkg-overview").NextUntil("#pkg-index")
	// ignore first import "pkgname" p tag
	overview := comments(sel)
	examples := examples(sel)
	url := sel.Find("code").First().Text()
	url = url[8 : len(url)-1]

	return &state{
		doc: document,
		pkg: doc.Package{
			URL:       url,
			Name:      name,
			Overview:  overview[1:],
			Examples:  examples,
			Functions: map[string]doc.Function{},
			Types:     map[string]doc.Type{},
		},
		useCase: useCase,
	}
}

func (s *state) newError(sel *goquery.Selection, msg string) error {
	return ParseError{sel, msg}
}

func (s *state) function(sel *goquery.Selection) error {
	next := sel.NextUntil("h3, h4")

	name, ok := sel.Attr("id")
	if !ok {
		return s.newError(sel, "could not get id")
	}
	signature := next.First().Text()

	f := doc.Function{
		Name:      name,
		Signature: signature,
		Comment:   comments(next),
		Examples:  examples(next),
	}

	if !s.useCase {
		name = strings.ToLower(name)
	}

	s.pkg.Functions[name] = f
	if s.current != nil {
		s.current.TypeFunctions[name] = f
	}
	return nil
}

func (s *state) typ(sel *goquery.Selection) error {
	if s.current != nil {
		name := s.current.Name
		if !s.useCase {
			name = strings.ToLower(name)
		}
		s.pkg.Types[name] = *s.current
	}

	next := sel.NextUntil("h3, h4")
	name, ok := sel.Attr("id")
	if !ok {
		return s.newError(sel, "could not get id")
	}
	signature := next.First().Text()

	t := doc.Type{
		Name:          name,
		Signature:     signature,
		Comment:       comments(next),
		Examples:      examples(next),
		TypeFunctions: map[string]doc.Function{},
		Methods:       map[string]doc.Method{},
	}

	s.current = &t
	return nil
}

func (s *state) method(sel *goquery.Selection) error {
	if s.current == nil {
		return s.newError(sel, "could not get method type")
	}

	next := sel.NextUntil("h3, h4")
	name, ok := sel.Attr("id")
	if !ok {
		return s.newError(sel, "could not get id")
	}
	split := strings.SplitN(name, ".", 2)
	name = split[len(split)-1]

	signature := next.First().Text()

	m := doc.Method{
		For: split[0],
		Function: doc.Function{
			Name:      name,
			Signature: signature,
			Comment:   comments(next),
			Examples:  examples(next),
		},
	}

	if !s.useCase {
		name = strings.ToLower(name)
	}

	s.current.Methods[name] = m
	return nil
}

func comments(sel *goquery.Selection) doc.Comment {
	nodes := sel.Filter("p, pre").Nodes
	comments := make(doc.Comment, 0, len(nodes))

	for _, node := range nodes {
		text := node.FirstChild.Data
		switch node.Data {
		case "p":
			var text string

			for child := node.FirstChild; child != nil; child = child.NextSibling {
				switch child.Type {
				case html.TextNode:
					text += child.Data
				case html.ElementNode:
					text += child.FirstChild.Data
				}
			}

			f := strings.Fields(text)
			comments = append(comments, doc.Paragraph(strings.Join(f, " ")))
		case "pre":
			comments = append(comments, doc.Pre(text))
		}
	}
	return comments
}

func examples(sel *goquery.Selection) []doc.Example {
	sel = sel.Find(".panel")
	examples := make([]doc.Example, 0, len(sel.Nodes))
	sel.Each(func(_ int, s *goquery.Selection) {
		// typically "Example¶"
		name := s.Find("summary").Text()

		pre := s.Find("pre")
		code, output := pre.First().Text(), pre.Last().Text()
		if code == output {
			output = ""
		}

		examples = append(examples, doc.Example{
			Name:   strings.TrimSuffix(name, "¶"),
			Code:   code,
			Output: output,
		})
	})
	return examples
}
