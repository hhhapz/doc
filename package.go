package doc

import (
	"html"
	"strings"
)

type Package struct {
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Overview Comment   `json:"overview"`
	Examples []Example `json:"examples"`

	Functions map[string]Function `json:"functions"`
	Types     map[string]Type     `json:"types"`
}

type Type struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Signature string    `json:"signature"`
	Comment   Comment   `json:"comment"`
	Examples  []Example `json:"examples"`

	TypeFunctions map[string]Function `json:"type_functions"`
	Methods       map[string]Method   `json:"methods"`
}

type Function struct {
	Name      string    `json:"name"`
	Signature string    `json:"signature"`
	Comment   Comment   `json:"comment"`
	Examples  []Example `json:"examples"`
}

type Method struct {
	For string `json:"for"`
	Function
}

type Note interface {
	Text() string
	HTML() string
	Markdown() string
}

var (
	_ Note = Comment(nil)
	_ Note = Paragraph("")
	_ Note = Pre("")
)

type Comment []Note

func (c Comment) Text() string {
	var s []string
	for _, n := range c {
		s = append(s, n.Text())
	}
	return strings.Join(s, "\n")
}

func (c Comment) HTML() string {
	var s []string
	for _, n := range c {
		s = append(s, n.HTML())
	}
	return strings.Join(s, "\n")
}

func (c Comment) Markdown() string {
	var s []string
	for _, n := range c {
		s = append(s, n.Markdown())
	}
	return strings.Join(s, "\n")
}

type Paragraph string

func (p Paragraph) Text() string {
	return string(p)
}

func (p Paragraph) HTML() string {
	return "<p>" + html.EscapeString(string(p)) + "</p>"
}

func (p Paragraph) Markdown() string {
	return string(p)
}

type Pre string

func (pre Pre) Text() string {
	s := strings.Split(string(pre), "\n")
	return strings.Join(s, "\n    ")
}

func (pre Pre) HTML() string {
	return "<pre>" + html.EscapeString(string(pre)) + "</pre>"
}

func (pre Pre) Markdown() string {
	return "```\n" + string(pre) + "\n```"
}

type Example struct {
	Name   string
	Code   string
	Output string
}