package doc

import "context"

type Searcher interface {
	// Search will find a package with the module name.
	Search(ctx context.Context, module string) (Package, error)
}
