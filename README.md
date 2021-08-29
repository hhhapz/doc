# doc

[![version][goversion]][go-dev]
[![version][pkgbadge]][pkglink]

This module provides an API to programatically search the documentation of Go
modules.

## Usage

To import and use, `go get github.com/hhhapz/doc`

```go
s := doc.New(http.DefaultClient, godocs.Parser) // or pkgsite.Parser
pkg, err := s.Search(context.TODO(), "bytes")

// use pkg
```

### Caching packages

The doc package also has a basic caching implementation that stores results in
an in-memory map.

```go
s := doc.New(http.DefaultClient, godocs.Parser) // or pkgsite.Parser
cs := doc.WithCache(s)

pkg, err := cs.Search(context.TODO(), "bytes")

// Cached results
pkg, err := cs.Search(context.TODO(), "bytes")
```

---

This package relies on [https://godocs.io][godocs].
It is planned to add a parser for [pkgsite][pkgsite] as well.

<!-- -->
[godocs]: https://godocs.io
[go-dev]: https://go.dev
[pkgsite]: https://pkg.go.dev
[pkglink]: https://pkg.go.dev/badge/github.com/hhhapz/doc.svg
<!-- -->
[goversion]: https://img.shields.io/github/go-mod/go-version/hhhapz/doc?color=%23007D9C&label=Go&style=flat
[pkgbadge]: https://pkg.go.dev/badge/github.com/hhhapz/doc.svg
