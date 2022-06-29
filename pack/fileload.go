package pack

import (
	"os"
	"sync"

	"github.com/hsiaosiyuan0/mole/util"
)

type FileLoadResult struct {
	raw []byte
	err error
}

type FileLoader struct {
	// path => code
	cache *util.LruCache[string, []byte]

	loading map[string][]chan *FileLoadResult
	lock    sync.Mutex
}

func NewFileLoader(cap int, clear int) *FileLoader {
	return &FileLoader{
		cache:   util.NewLruCache[string, []byte](cap, clear),
		loading: map[string][]chan *FileLoadResult{},
		lock:    sync.Mutex{},
	}
}

func (f *FileLoader) Load(file string) (interface{}, error) {
	if f.cache.HasKey(file) {
		return f.cache.Get(file), nil
	}

	// subscribe
	f.lock.Lock()
	if queue, ok := f.loading[file]; ok {
		c := make(chan *FileLoadResult)
		f.loading[file] = append(queue, c)
		f.lock.Unlock()
		return c, nil
	}

	// set loading
	f.loading[file] = []chan *FileLoadResult{}
	f.lock.Unlock()

	code, err := os.ReadFile(file)
	var r *FileLoadResult
	if err != nil {
		r = &FileLoadResult{nil, err}
	} else {
		r = &FileLoadResult{code, nil}
	}

	// notify others
	f.lock.Lock()
	for _, c := range f.loading[file] {
		c <- r
	}
	delete(f.loading, file)
	f.lock.Unlock()

	if err != nil {
		return nil, err
	}

	f.cache.Set(file, code)
	return code, nil
}
