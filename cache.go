package doc

import (
	"context"
	"sync"
	"time"
)

type CachedSearcher struct {
	searcher Searcher

	mu    sync.Mutex
	cache map[string]*CachedPackage
}

type CachedPackage struct {
	Package
	Created time.Time
	Updated time.Time
}

func (cs *CachedSearcher) Search(ctx context.Context, module string) (Package, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if pkg, ok := cs.cache[module]; ok {
		pkg.Updated = time.Now()
		return pkg.Package, nil
	}

	pkg, err := cs.searcher.Search(ctx, module)
	if err != nil {
		return Package{}, err
	}

	cs.cache[module] = &CachedPackage{
		Package: pkg,
		Created: time.Now(),
		Updated: time.Now(),
	}
	return pkg, nil
}

func (cs *CachedSearcher) WithCache(f func(cache map[string]*CachedPackage)) {
	cs.mu.Lock()
	f(cs.cache)
	cs.mu.Unlock()
}
