package main

import (
	"sort"

	"github.com/hhhapz/doc"
)

const sdTpl = `
{{- $pkg := . -}}
• Overview
• Constants
• Variables
• Functions
{{- range $f := sortedRange .Functions }}
{{- $func := index $pkg.Functions $f }}
 • {{ index $func.Name }}
{{- end }}
• Types
{{- range $t := sortedRange .Types }}
{{- $typ := index $pkg.Types $t }}
 • {{ $typ.Name }}
{{- range $tf := sortedRange $typ.TypeFunctions }}
{{- $tfunc := index $typ.TypeFunctions $tf }}
  • {{ $tfunc.Name }}
{{- end }}
{{- range $m := sortedRange $typ.Methods }}
{{- $method := index $typ.Methods $m }}
  • {{ $method.Name }}
{{- end }}
{{- end }}`

const pkgTpl = `
{{- $pkg := . -}}
#  Package {{ .Name }}
---
{{ .Overview.Markdown }}

---

# Constants

{{ range .Constants }}
` + "```go" + `
{{ .Signature }}
` + "```" + `

{{ .Comment.Markdown }}

---
{{ else }}
This section is empty

---
{{ end }}

# Variables

{{ range .Variables }}
` + "```go" + `
{{ .Signature }}
` + "```" + `

{{ .Comment.Markdown }}

---
{{ else }}
This section is empty

---
{{ end }}

# Functions

{{ range $f := sortedRange .Functions }}
{{ $data := index $pkg.Functions $f }}
## func {{ $data.Name }}
` + "```go" + `
{{ $data.Signature }}
` + "```" + `

{{ $data.Comment.Markdown }}

---
{{ else }}
This section is empty

---
{{ end }}

# Types

{{ range $t := sortedRange .Types}}
{{ $data := index $pkg.Types $t }}
## type {{ $data.Name }}
` + "```go" + `
{{ $data.Signature }}
` + "```" + `

{{ $data.Comment.Markdown }}

---

{{- range $tf := sortedRange $data.TypeFunctions}}
{{ $typeFunc := index $data.TypeFunctions $tf }}
## func {{ $typeFunc.Name }}
` + "```go" + `
{{ $typeFunc.Signature }}
` + "```" + `

{{ $typeFunc.Comment.Markdown }}

---
{{ end }}


{{- range $m := sortedRange $data.Methods}}
{{ $method := index $data.Methods $m }}
## func {{ $method.Name }}
` + "```go" + `
{{ $method.Signature }}
` + "```" + `

{{ $method.Comment.Markdown }}

---
{{ end }}

{{ else }}
This section is empty

---
{{ end }}
`

func sortedRange(m any) []string {
	var keys []string
	switch v := m.(type) {
	case map[string]doc.Variable:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	case map[string]doc.Function:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	case map[string]doc.Type:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	case map[string]doc.Method:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
