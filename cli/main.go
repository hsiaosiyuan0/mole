package main

import (
	"flag"
	"os"
)

type Options struct {
	ast  bool
	file string

	packAna bool
	dir     string
	cfg     string
}

func newOptions() *Options {
	opts := &Options{}

	flag.BoolVar(&opts.ast, "ast", false, "print AST of the target file")
	flag.StringVar(&opts.file, "file", "", "print AST of the target file")

	flag.BoolVar(&opts.packAna, "pkg-ana", false, "analyze the package information")
	flag.StringVar(&opts.dir, "dir", "", "the project directory")
	flag.StringVar(&opts.cfg, "cfg", "", "the config file")

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
	cmds := &[]SubCommand{&AstInspector{}, &PkgAnalysis{}}
	for _, cmd := range *cmds {
		if cmd.Process(opts) {
			return
		}
	}
	flag.Usage()
}
