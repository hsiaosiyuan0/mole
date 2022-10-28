package main

import (
	"flag"
	"os"
)

type Options struct {
	ast  bool
	file string

	dir string
	cfg string
	out string

	perf bool
}

func newOptions() *Options {
	opts := &Options{}

	flag.BoolVar(&opts.ast, "ast", false, "print AST of the target file")
	flag.StringVar(&opts.file, "file", "", "print AST of the target file")

	flag.StringVar(&opts.dir, "dir", "", "the project directory")
	flag.StringVar(&opts.cfg, "cfg", "", "the config file")
	flag.StringVar(&opts.out, "out", "", "the output file")

	flag.BoolVar(&opts.perf, "perf", false, "gen the pprof file")

	flag.Parse()

	if opts.dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		opts.dir = cwd
	}

	return opts
}

type SubCommand interface {
	Process(*Options) bool
}

func main() {
	opts := newOptions()
	cmds := &[]SubCommand{&AstInspector{}}
	for _, cmd := range *cmds {
		if cmd.Process(opts) {
			return
		}
	}
	flag.Usage()
}
