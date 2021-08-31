package doc

import (
	"context"
	"net/http"
	"sync"
)

type Searcher interface {
	// Search will find a package with the module name.
	Search(ctx context.Context, module string) (Package, error)
}

type configurer interface {
	withAgent(string)
	maintainCase()
}

func UserAgent(agent string) func(configurer) {
	return func(c configurer) {
		c.withAgent(agent)
	}
}

func MaintainCase() func(configurer) {
	return func(c configurer) {
		c.maintainCase()
	}
}

func New(client *http.Client, parser Parser, opts ...func(configurer)) *HTTPSearcher {
	s := &HTTPSearcher{
		client:   client,
		parser:   parser,
		agent:    "Doc (https://github.com/hhhapz/doc)",
		withCase: false,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithCache(s Searcher) *CachedSearcher {
	return &CachedSearcher{
		searcher: s,
		mu:       sync.Mutex{},
		cache:    map[string]*CachedPackage{},
	}
}
