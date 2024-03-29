package util

import (
	"container/list"
	"os"
	"path/filepath"
	"sync"
)

type DirWalkerHandle = func(string, bool, *DirWalker)

// by the benchmark, there is no execution performance benefit from the works pool,
// however the works pool is still useful since it can stop the routine more precisely
type DirWalker struct {
	Dir        string
	Concurrent int

	handle DirWalkerHandle

	dirs     *list.List
	dirsLock sync.Mutex

	newJob chan bool

	wg    sync.WaitGroup
	wgFin chan bool

	stop chan error
	fin  chan bool
	err  error
}

func NewDirWalker(dir string, concurrent int, handle DirWalkerHandle) *DirWalker {
	if concurrent == 0 {
		concurrent = 16
	}

	w := &DirWalker{
		Dir:        dir,
		Concurrent: concurrent,
		handle:     handle,

		dirs:     list.New(),
		dirsLock: sync.Mutex{},

		newJob: make(chan bool),

		wg:    sync.WaitGroup{},
		wgFin: make(chan bool),

		stop: make(chan error),
		fin:  make(chan bool),
	}

	return w.initWorkers()
}

func (w *DirWalker) shift() string {
	if w.dirs.Len() == 0 {
		return ""
	}
	w.dirsLock.Lock()
	defer w.dirsLock.Unlock()

	d := w.dirs.Front()
	w.dirs.Remove(d)
	return d.Value.(string)
}

func (w *DirWalker) push(file string) {
	w.dirsLock.Lock()
	defer w.dirsLock.Unlock()

	w.dirs.PushBack(file)
}

func (w *DirWalker) walk() {
loop:
	for {
		select {
		case <-w.newJob:
			dir := w.shift()
			if dir == "" {
				continue
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				w.stop <- err
				return
			}

			// process the handle of directory synchronously to keep the lexical order
			// within the directory and its children files
			w.handle(dir, true, w)

			for _, file := range files {
				pth := filepath.Join(dir, file.Name())
				if file.IsDir() {
					w.wg.Add(1)
					w.push(pth)
					go func() { w.newJob <- true }()
				} else {
					w.handle(pth, false, w)
				}
			}

			w.wg.Done()
		case <-w.fin:
			break loop
		}
	}
}

func (w *DirWalker) initWorkers() *DirWalker {
	for i := 0; i < w.Concurrent; i++ {
		go w.walk()

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
	w.wg.Add(1)
	w.push(w.Dir)
	w.newJob <- true

	go func() {
		w.wg.Wait()
		w.wgFin <- true
	}()

loop:
	for {
		select {
		case <-w.wgFin:
			break loop
		case err := <-w.stop:
			w.err = err
			break loop
		}
	}

	w.fin <- true
}
