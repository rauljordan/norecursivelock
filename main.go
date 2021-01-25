package main

import (
	"github.com/rauljordan/nolock/readlock"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(readlock.Analyzer)
}
