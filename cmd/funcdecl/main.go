package main

import (
	"github.com/ichiban/funcdecl"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(funcdecl.Analyzer)
}
