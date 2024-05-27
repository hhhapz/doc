package doc

import (
	"context"
	"sync"
	"time"
)

type cachedSearcher struct {
	Searcher

	mu    sync.RWMutex
	cache map[string]*CachedPackage
}

type CachedPackage struct {
	Package
	Created time.Time
	Updated time.Time
}

func (c *cachedSearcher) Search(ctx context.Context, module string) (Package, error) {
	c.mu.RLock()
	cPkg, ok := c.cache[module]
	c.mu.RUnlock()
	if ok {
		cPkg.Updated = time.Now()
		return cPkg.Package, nil
	}

	pkg, err := c.Searcher.Search(ctx, module)
	if err != nil {
		return Package{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[module] = &CachedPackage{
		Package: pkg,
		Created: time.Now(),
		Updated: time.Now(),
	}
	return pkg, nil
}

func (c *cachedSearcher) WithCache(f func(cache map[string]*CachedPackage)) {
	c.mu.Lock()
	f(c.cache)
	c.mu.Unlock()
}
