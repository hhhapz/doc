package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hhhapz/doc"
	"github.com/hhhapz/doc/pkgsite"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if flag.Arg(0) == "" {
		return fmt.Errorf("usage: %s <package name>", os.Args[0])
	}
	s := doc.NewSearcher(pkgsite.Parser)
	pkg, err := s.Search(ctx, flag.Arg(0))
	if err != nil {
		return fmt.Errorf("could not fetch package: %w", err)
	}

	model, err := newModel(pkg, *debug)
	if err != nil {
		return fmt.Errorf("could not initialize: %v", err)
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		return fmt.Errorf("could not run: %v", err)
	}
	return nil
}
