package godocs_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hhhapz/doc"
	"github.com/hhhapz/doc/godocs"
)

func TestGodocs(t *testing.T) {
	ctx := context.Background()

	s := doc.NewSearcher(godocs.Parser)
	pkg, err := s.Search(ctx, "net")
	if err != nil {
		t.Errorf("could not fetch stdlib package: %v", err)
		return
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(pkg)
}
