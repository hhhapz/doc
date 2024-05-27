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

func NewSearcher(parser Parser, opts ...SearchOption) Searcher {
	s := &httpSearcher{
		client:   http.DefaultClient,
		parser:   parser,
		agent:    "Doc (https://github.com/hhhapz/doc)",
		withCase: false,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type CachedSearcher interface {
	Searcher
	// WithCache gives access to modify and update the contents of the internal cache.
	WithCache(func(cache map[string]*CachedPackage))
}

func NewCachedSearcher(parser Parser, opts ...SearchOption) CachedSearcher {
	s := NewSearcher(parser, opts...)
	return &cachedSearcher{
		Searcher: s,
		mu:       sync.RWMutex{},
		cache:    map[string]*CachedPackage{},
	}
}

type SearchOption = func(s *httpSearcher)

func WithClient(client *http.Client) SearchOption {
	return func(s *httpSearcher) {
		s.client = client
	}
}

func UserAgent(agent string) SearchOption {
	return func(s *httpSearcher) {
		s.agent = agent
	}
}

func MaintainCase() SearchOption {
	return func(s *httpSearcher) {
		s.withCase = true
	}
}
