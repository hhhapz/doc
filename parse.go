package doc

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	pkg     Package
	current *Type
}

func NewState(doc *goquery.Document) *state {
	name := doc.Find("#pkg-overview").Text()
	name = strings.TrimPrefix(name, "package ")

	sel := doc.Find("#pkg-overview").NextUntil("#pkg-index")
	// ignore first import "pkgname" p tag
	overview := comments(sel)
	examples := examples(sel)
	url := sel.Find("code").First().Text()
	url = url[8 : len(url)-1]

	return &state{
		doc: doc,
		pkg: Package{
			URL:       url,
			Name:      name,
			Overview:  overview[1:],
			Examples:  examples,
			Functions: map[string]Function{},
			Types:     map[string]Type{},
		},
	}
}

func (s *state) NewError(sel *goquery.Selection, msg string) error {
	return ParseError{sel, msg}
}

func (s *state) function(sel *goquery.Selection) error {
	next := sel.NextUntil("h3, h4")

	name, ok := sel.Attr("id")
	if !ok {
		return s.NewError(sel, "could not get id")
	}
	signature := next.First().Text()

	f := Function{
		Name:      name,
		Signature: signature,
		Comment:   comments(next),
		Examples:  examples(next),
	}

	s.pkg.Functions[name] = f
	if s.current != nil {
		s.current.TypeFunctions[name] = f
	}
	return nil
}

func (s *state) typ(sel *goquery.Selection) error {
	if s.current != nil {
		s.pkg.Types[s.current.Name] = *s.current
	}
	next := sel.NextUntil("h3, h4")

	name, ok := sel.Attr("id")
	if !ok {
		return s.NewError(sel, "could not get id")
	}
	signature := next.First().Text()

	t := Type{
		Name:          name,
		Signature:     signature,
		Comment:       comments(next),
		Examples:      examples(next),
		TypeFunctions: map[string]Function{},
		Methods:       map[string]Method{},
	}

	s.current = &t
	return nil
}

func (s *state) method(sel *goquery.Selection) error {
	if s.current == nil {
		return s.NewError(sel, "could not get method type")
	}

	next := sel.NextUntil("h3, h4")
	name, ok := sel.Attr("id")
	if !ok {
		return s.NewError(sel, "could not get id")
	}
	split := strings.SplitN(name, ".", 2)
	name = split[len(split)-1]

	signature := next.First().Text()

	m := Method{
		For: split[0],
		Function: Function{
			Name:      name,
			Signature: signature,
			Comment:   comments(next),
			Examples:  examples(next),
		},
	}
	s.current.Methods[name] = m
	return nil
}

func comments(sel *goquery.Selection) Comment {
	nodes := sel.Filter("p, pre").Nodes
	comments := make(Comment, 0, len(nodes))

	for _, node := range nodes {
		text := node.FirstChild.Data
		switch node.Data {
		case "p":
			f := strings.Fields(text)
			text = strings.Join(f, " ")
			comments = append(comments, Paragraph(text))
		case "pre":
			comments = append(comments, Pre(text))
		}
	}
	return comments
}

func examples(sel *goquery.Selection) []Example {
	sel = sel.Find(".panel")
	examples := make([]Example, 0, len(sel.Nodes))
	sel.Each(func(_ int, s *goquery.Selection) {
		// typically "Example¶"
		name := s.Find("summary").Text()

		pre := s.Find("pre")
		code, output := pre.First().Text(), pre.Last().Text()
		if code == output {
			output = ""
		}

		examples = append(examples, Example{
			Name:   strings.TrimSuffix(name, "¶"),
			Code:   code,
			Output: output,
		})
	})
	return examples
}
