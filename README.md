# doc

This module provides an API to programatically search the documentation of Go
modules.

To import and use, `go get github.com/hhhapz/doc`

```go
s := doc.New(http.DefaultClient, godocs.Parser) // or pkgsite.Parser
pkg, err := s.Search(context.TODO(), "bytes")

// use pkg
```

---

This package relies on [https://godocs.io](https://godocs.io).
It is planned to add a parser for [pkgsite](https://pkg.go.dev) as well.
