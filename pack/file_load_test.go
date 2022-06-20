package pack

import (
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"testing"
)

func (f *FileLoader) testLoadNoChain(file string) ([]byte, error) {
	if f.cache.HasKey(file) {
		return f.cache.Get(file), nil
	}

	code, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	f.cache.Set(file, code)
	return code, nil
}

func BenchmarkFileLoaderNoCache(b *testing.B) {
	// create sample files
	cnt := 5000
	files := []string{}

	for i := 0; i < cnt; i++ {
		file, err := os.CreateTemp("", "BenchmarkFileLoader-*")
		if err != nil {
			b.Fatal(err)
		}
		files = append(files, file.Name())
		ioutil.WriteFile(file.Name(), make([]byte, 4096), 0644)
	}
	defer func() {
		for _, file := range files {
			os.Remove(file)
		}
	}()

	// init workers
	wg := sync.WaitGroup{}
	theFiles := make(chan string)

	workerCnt := 10
	var bytesNum uint64
	for i := 0; i < workerCnt; i++ {
		go func() {
			for {
				select {
				case file := <-theFiles:
					f, err := os.ReadFile(file)
					if err != nil {
						b.Log(err)
					}
					atomic.AddUint64(&bytesNum, uint64(len(f)))
					wg.Done()
				}
			}
		}()
	}

	b.ResetTimer()

	rCnt := 3 * cnt
	for i := 0; i < rCnt; i++ {
		wg.Add(1)
		file := files[rand.Intn(cnt)]
		theFiles <- file
	}

	wg.Wait()

	if bytesNum != uint64(rCnt*4096) {
		b.Fatal("failed")
	}
}

func BenchmarkFileLoaderNoLoadingBarrier(b *testing.B) {
	// create sample files
	cnt := 5000
	files := []string{}

	for i := 0; i < cnt; i++ {
		file, err := os.CreateTemp("", "BenchmarkFileLoader-*")
		if err != nil {
			b.Fatal(err)
		}
		files = append(files, file.Name())
		ioutil.WriteFile(file.Name(), make([]byte, 4096), 0644)
	}
	defer func() {
		for _, file := range files {
			os.Remove(file)
		}
	}()

	// init workers
	wg := sync.WaitGroup{}
	theFiles := make(chan string)
	loader := NewFileLoader(1024, 10)

	workerCnt := 10
	var bytesNum uint64
	for i := 0; i < workerCnt; i++ {
		go func() {
			for {
				select {
				case file := <-theFiles:
					f, err := loader.testLoadNoChain(file)
					if err != nil {
						b.Log(err)
					}
					atomic.AddUint64(&bytesNum, uint64(len(f)))
					wg.Done()
				}
			}
		}()
	}

	b.ResetTimer()

	rCnt := 3 * cnt
	for i := 0; i < rCnt; i++ {
		wg.Add(1)
		file := files[rand.Intn(cnt)]
		theFiles <- file
	}

	wg.Wait()

	if bytesNum != uint64(rCnt*4096) {
		b.Fatal("failed")
	}
}

func BenchmarkFileLoader(b *testing.B) {
	// create sample files
	cnt := 5000
	files := []string{}

	for i := 0; i < cnt; i++ {
		file, err := os.CreateTemp("", "BenchmarkFileLoader-*")
		if err != nil {
			b.Fatal(err)
		}
		files = append(files, file.Name())
		ioutil.WriteFile(file.Name(), make([]byte, 4096), 0644)
	}
	defer func() {
		for _, file := range files {
			os.Remove(file)
		}
	}()

	// init workers
	wg := sync.WaitGroup{}
	theFiles := make(chan string)
	loader := NewFileLoader(1024, 10)

	workerCnt := 10
	var bytesNum uint64
	for i := 0; i < workerCnt; i++ {
		go func() {
			for {
				select {
				case file := <-theFiles:
					f, err := loader.Load(file)
					if err != nil {
						b.Log(err)
					}
					switch fv := f.(type) {
					case []byte: // done
						atomic.AddUint64(&bytesNum, uint64(len(fv)))
						wg.Done()
					case chan *FileLoadResult:
						f := <-fv // wait
						atomic.AddUint64(&bytesNum, uint64(len(f.raw)))
						wg.Done()
					}
				}
			}
		}()
	}

	b.ResetTimer()

	rCnt := 3 * cnt
	for i := 0; i < rCnt; i++ {
		wg.Add(1)
		file := files[rand.Intn(cnt)]
		theFiles <- file
	}

	wg.Wait()

	if bytesNum != uint64(rCnt*4096) {
		b.Fatal("failed")
	}
}
