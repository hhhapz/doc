package doc

import (
	"context"
	"net/http"
	"time"
)

type CachedSearcher interface {
	Searcher

	Add(module string, pkg Package)
	All() (packages []Package)
	Cached(module string) (exists bool)
	Clear(module string) (exists bool)
	ClearAll() (amount int)
	Prune(before time.Time) (amount int)
}

type Searcher interface {
	// Search will find a package with the module name.
	Search(ctx context.Context, module string) (Package, error)
}

type Configurer interface {
	SetAgent(string)
}

func UserAgent(agent string) func(Configurer) {
	return func(c Configurer) {
		c.SetAgent(agent)
	}
}

func New(client *http.Client, parser Parser, opts ...func(Configurer)) Searcher {
	s := &httpSearcher{
		Client: client,
		Parser: parser,
		Agent:  "Doc (https://github.com/hhhapz/doc)",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func NewWithCache(client *http.Client, parser Parser) CachedSearcher {
	return nil
}
