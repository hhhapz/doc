package pkgsite_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hhhapz/doc"
	"github.com/hhhapz/doc/pkgsite"
)

func TestPkgSite(t *testing.T) {
	ctx := context.Background()
	s := doc.NewSearcher(pkgsite.Parser)
	// testPackage(t, ctx, s, "database/sql")
	// testPackage(t, ctx, s, "github.com/diamondburned/arikawa/v3")
	// testPackage(t, ctx, s, "github.com/hhhapz/diffgen")
	testPackage(t, ctx, s, "net")
}

func testPackage(t *testing.T, ctx context.Context, s doc.Searcher, pkgName string) {
	pkg, err := s.Search(ctx, pkgName)
	if err != nil {
		t.Errorf("could not fetch stdlib package: %v", err)
		return
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(pkg)
}
