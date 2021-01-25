package main

import (
	"github.com/rauljordan/nolock/deadreadlock"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(deadreadlock.Analyzer)
}
