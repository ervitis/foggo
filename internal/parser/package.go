package parser

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type PackageInfo struct {
	Name     string
	AstFiles []*ast.File
	Paths    []string
}

func ParsePackageInfo(p string) (*PackageInfo, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedFiles,
		Tests: false,
	}

	packagesList, err := packages.Load(cfg, p)
	if err != nil {
		return nil, err
	}

	if len(packagesList) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(packagesList))
	}

	pl := packagesList[0]
	return &PackageInfo{
		Name:     pl.Name,
		AstFiles: pl.Syntax,
		Paths:    pl.GoFiles,
	}, nil
}
