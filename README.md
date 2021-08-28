# doc

This module provides an API to programatically search the documentation of Go
modules.

To import and use, `go get github.com/hhhapz/doc`

```go
s := doc.HTTPSearcher{}
pkg, err := s.Search("bytes")

// use pkg
```

---

This package relies on [https://godoc.io](https://godoc.io).
It is planned to add a parser for [pkgsite](https://pkg.go.dev) as well.
