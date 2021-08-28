package doc

import "context"

type Searcher interface {
	// Search will find a package with the module name.
	Search(module string) (Package, error)

	// SearchContext will find a package with the module name.
	SearchContext(ctx context.Context, module string) (Package, error)
}
