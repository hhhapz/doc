package doc

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	// Base is the base path to godocs.io
	Base = "https://godocs.io/"
	// DefaultuserAgent is the the default user agent when not provided.
	DefaultUserAgent = "Doc (https://github.com/hhhapz/doc)"
	// Selectors matches methods, types, and functions on the documentation
	// page.
	Selectors = `
[data-kind="function"],
[data-kind="type"],
[data-kind="method"]:not([class*="decl"])
`
)

// InvalidStatusError indicates that the request to the godocs.io was not
// successful. The value is the status that was returned from the page instead.
type InvalidStatusError int

// Error satisfies the error interface.
func (err InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid response status: %d", err)
}

// HTTPSearcher provides an interface to search the godocs package module page.
// It implements the Searcher interface. The zero value is ready to use.
//
// HTTPSearcher does not cache results and will do the request every time, even
// if provided the same module name. If caching is required, the CachedSearcher
// type.
type HTTPSearcher struct {
	Client *http.Client
	Agent  string
}

// HTTPSearcher implements the Searcher interface.
var _ Searcher = HTTPSearcher{}

// Search searches godocs for the provided module.
//
// Search calls SearchContext with context.Background().
func (h HTTPSearcher) Search(module string) (Package, error) {
	return h.SearchContext(context.Background(), module)
}

// Search searches godocs for the provided module.
//
// SearchContext is the main workhorse for querying and parsing the http
// response. The implementation for parsing the document can be found in
// parse.go
//
// If the page does not respond with a 200 status code, a InvalidStatusError is
// returned. If the page could not be parsed by GoQuery, the error will be of
// type Otherwise, issues while parsing the document will of type ParseError,
// and will contain the selector being parsed, for more context.
func (h HTTPSearcher) SearchContext(ctx context.Context, module string) (Package, error) {
	body, err := h.request(ctx, module)
	if err != nil {
		return Package{}, err
	}

	doc, err := goquery.NewDocumentFromReader(body)
	_ = body.Close()
	if err != nil {
		return Package{}, err
	}

	s := newState(doc)

	doc.Find(`[data-kind="function"], [data-kind="type"], [data-kind="method"]:not([class*="decl"])`).EachWithBreak(func(_ int, sel *goquery.Selection) bool {
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
		return Package{}, err
	}
	return s.pkg, nil
}

// request is a helper function to do the http request and return the body.
func (h HTTPSearcher) request(ctx context.Context, module string) (io.ReadCloser, error) {
	r, err := http.NewRequestWithContext(ctx, "GET", Base+module, http.NoBody)
	if err != nil {
		return nil, err
	}

	if h.Agent == "" {
		h.Agent = DefaultUserAgent
	}
	r.Header.Add("User-Agent", h.Agent)

	if h.Client == nil {
		h.Client = http.DefaultClient
	}
	resp, err := h.Client.Do(r)
	if err != nil {
		return nil, err
	}

	if c := resp.StatusCode; c != 200 {
		return nil, InvalidStatusError(c)
	}

	return resp.Body, nil
}
