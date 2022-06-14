package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestWalkDir(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test")

	var lock sync.Mutex
	files := []string{}
	w := NewDirWalker(dir, 0, func(f string, dir bool, dw *DirWalker) {
		lock.Lock()
		defer lock.Unlock()

		files = append(files, f)
	})

	w.Walk()

	AssertEqual(t, 23, len(files), "should be ok")
}

func TestWalkDirStop(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test")

	w := NewDirWalker(dir, 0, func(f string, dir bool, dw *DirWalker) {
		dw.Stop(errors.New("stopped"))
	})

	w.Walk()

	AssertEqual(t, "stopped", w.Err().Error(), "should be ok")
}

func BenchmarkWalkOneWorker(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}

	cwd, err = filepath.Abs(path.Join(cwd, ".."))
	if err != nil {
		b.Fatal(err)
	}

	var lock sync.Mutex
	files := []string{}
	w := NewDirWalker(cwd, 1, func(f string, dir bool, dw *DirWalker) {
		lock.Lock()
		defer lock.Unlock()

		files = append(files, f)
	})
	w.Walk()

	fmt.Printf("files count: %d\n", len(files))
}

func BenchmarkWalkMultiWorkers(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}

	cwd, err = filepath.Abs(path.Join(cwd, ".."))
	if err != nil {
		b.Fatal(err)
	}

	var lock sync.Mutex
	files := []string{}
	w := NewDirWalker(cwd, 0, func(f string, dir bool, dw *DirWalker) {
		lock.Lock()
		defer lock.Unlock()

		files = append(files, f)
	})
	w.Walk()

	fmt.Printf("files count: %d\n", len(files))
}
