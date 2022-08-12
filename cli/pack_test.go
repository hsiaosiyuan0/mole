package main

import (
	"testing"
)

func TestPkgAna(t *testing.T) {
	ana := &PkgAnalysis{}
	opts := &Options{
		packAna: true,
		dir:     "/Users/hsiao/Developer/work/mole-tests/samples-data/web/st-year2021",
	}

	ana.Process(opts)
}
