package util

import (
	"os"
	"path"
	"runtime"
	"sync"
)

type DirWalkerHandle = func(string, bool, *DirWalker)

type DirWalker struct {
	Dir        string
	Concurrent int

	handle DirWalkerHandle

	dirs   chan string
	files  chan string
	wg     sync.WaitGroup
	wgChan chan bool

	stop chan error
	err  error
}

func NewDirWalker(dir string, concurrent int, handle DirWalkerHandle) *DirWalker {
	if concurrent == 0 {
		concurrent = runtime.NumCPU()
	}

	w := &DirWalker{
		Dir:        dir,
		Concurrent: concurrent,
		handle:     handle,

		dirs:   make(chan string, concurrent),
		files:  make(chan string, concurrent),
		wg:     sync.WaitGroup{},
		wgChan: make(chan bool),

		stop: make(chan error),
	}

	return w.initWorkers()
}

func (w *DirWalker) walk() {
	for dir := range w.dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			w.stop <- err
			return
		}

		// process the handle of directory synchronously to keep the lexical order
		// between the directory and its children files
		w.handle(dir, true, w)

		for _, file := range files {
			pth := path.Join(dir, file.Name())
			w.addFile(pth, file.IsDir())
		}

		w.wg.Done()
	}
}

func (w *DirWalker) addFile(file string, isDir bool) {
	w.wg.Add(1)

	go func(dir bool) {
		if dir {
			w.dirs <- file
		} else {
			w.files <- file
		}
	}(isDir)
}

func (w *DirWalker) doHandle() {
	for file := range w.files {
		w.handle(file, false, w)
		w.wg.Done()
	}
}

func (w *DirWalker) initWorkers() *DirWalker {
	for i := 0; i < w.Concurrent; i++ {
		go w.walk()
		go w.doHandle()
	}
	return w
}

func (w *DirWalker) Stop(err error) {
	w.stop <- err
}

func (w *DirWalker) Err() error {
	return w.err
}

func (w *DirWalker) Walk() {
	go func() {
		w.wg.Wait()
		w.wgChan <- true
	}()

	w.addFile(w.Dir, true)

loop:
	for {
		select {
		case <-w.wgChan:
			break loop
		case err := <-w.stop:
			w.err = err
			break loop
		}
	}
}
