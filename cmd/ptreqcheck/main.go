package main

import (
	"github.com/MortenLohne/ptreqcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(ptreqcheck.Analyzer)
}
