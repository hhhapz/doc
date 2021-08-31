# doc

[![version][goversion]][go-dev]
[![reference][pkgbadge]][pkglink]
[![tag][tagbadge]][pkglink]

This module provides an API to programatically search the documentation of Go
modules.

## Usage

To import and use, `go get github.com/hhhapz/doc`

```go
s := doc.New(http.DefaultClient, godocs.Parser, opts...) // or pkgsite.Parser
pkg, err := s.Search(context.TODO(), "bytes")

// use pkg
```

### Options (opts...)

Currently there are two available options:

#### `doc.MaintainCase()`

By default, the maps in the Package struct will have lower case keys:

- `Package.Functions`
- `Package.Types`
- `Package.Types.TypeFunctions`
- `Package.Types.Methods`

When enabling MaintainCase, the keys to all of these functions will be retained
to their true case.

#### `doc.UserAgent(string)`

UserAgent will allow you to change the UA agent for all requests to the package
sites. by default it will link to this repository.

---

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
[pkglink]: https://pkg.go.dev/github.com/hhhapz/doc
<!-- -->
[goversion]: https://img.shields.io/github/go-mod/go-version/hhhapz/doc?color=%23007D9C&label=Go&style=flat
[tagbadge]: https://img.shields.io/github/v/tag/hhhapz/doc?color=%23007d9c
[pkgbadge]: https://pkg.go.dev/badge/github.com/hhhapz/doc.svg
