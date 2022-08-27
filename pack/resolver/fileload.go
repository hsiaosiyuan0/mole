package resolver

import (
	"os"
	"sync"

	"github.com/hsiaosiyuan0/mole/util"
)

type FileLoadResult struct {
	Raw []byte
	Err error
}

type FileLoader struct {
	// path => code
	cache *util.LruCache[string, *FileLoadResult]

	// file => subscribers
	loading map[string][]chan *FileLoadResult
	lock    sync.Mutex
}

func NewFileLoader(cap int, clear int) *FileLoader {
	return &FileLoader{
		cache:   util.NewLruCache[string, *FileLoadResult](cap, clear),
		loading: map[string][]chan *FileLoadResult{},
		lock:    sync.Mutex{},
	}
}

func (fl *FileLoader) Load(file string) chan *FileLoadResult {
	if fl.cache.HasKey(file) {
		c := make(chan *FileLoadResult, 1)
		c <- fl.cache.Get(file)
		return c
	}

	// let the caller to subscribe the result
	fl.lock.Lock()
	queue := fl.loading[file]
	first := false
	if queue == nil {
		first = true
		queue = []chan *FileLoadResult{}
		fl.loading[file] = queue
	}

	c := make(chan *FileLoadResult)
	fl.loading[file] = append(queue, c)
	fl.lock.Unlock()

	// do the actual file load if its the first call
	if first {
		go func() {
			raw, err := os.ReadFile(file)
			r := &FileLoadResult{raw, err}
			fl.cache.Set(file, r)

			// notify the subscribers
			fl.lock.Lock()
			for _, c := range fl.loading[file] {
				c <- r
			}
			fl.loading[file] = nil
			fl.lock.Unlock()
		}()
	}

	return c
}
