package main

import (
	"testing"
)

func TestPkgAna(t *testing.T) {
	ana := &PkgAnalysis{}
	opts := &Options{
		packAna: true,
		dir:     "/Users/hsiao/Developer/work/mole-tests/samples-data/web/st-live-tanabat20220801",
	}

	ana.Process(opts)
}
