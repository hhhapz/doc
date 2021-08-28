package doc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	// DefaultuserAgent is the the default user agent when not provided.
	DefaultUserAgent = "Doc (https://github.com/hhhapz/doc)"
)

// Parser is the interface that package site parsers implement.
type Parser interface {
	URL(module string) (full string)
	Parse(document *goquery.Document) (Package, error)
}

// InvalidStatusError indicates that the request to the godocs.io was not
// successful. The value is the status that was returned from the page instead.
type InvalidStatusError int

// Error satisfies the error interface.
func (err InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid response status: %d", err)
}

var ErrNoParser = errors.New("parser not provided")

// HTTPSearcher provides an interface to search the godocs package module page.
// It implements the Searcher interface. A parser must be provided, such as
// pkgsite.Parser, or godoc.Parser.
//
// HTTPSearcher does not cache results and will do the request every time, even
// if provided the same module name. If caching is required, the CachedSearcher
// type.
type HTTPSearcher struct {
	Parser Parser
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
	if h.Parser == nil {
		return Package{}, ErrNoParser
	}

	body, err := h.request(ctx, module)
	if err != nil {
		return Package{}, err
	}
	defer body.Close()

	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return Package{}, err
	}
	return h.Parser.Parse(document)
}

// request is a helper function to do the http request and return the body.
func (h HTTPSearcher) request(ctx context.Context, module string) (io.ReadCloser, error) {
	url := h.Parser.URL(module)
	r, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
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
