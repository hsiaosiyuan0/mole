package plugin

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"github.com/hsiaosiyuan0/mole/util"
)

func Resolve(cwd, name string) (*plugin.Plugin, error) {
	sep := string(filepath.Separator)
	if !strings.HasPrefix(cwd, sep) {
		return nil, errors.New(fmt.Sprintf("absolute path is required: %s", cwd))
	}

	p := path.Join(cwd, "node_modules", name, "build")

	if !util.FileExist(p) {
		if p == sep {
			return nil, os.ErrNotExist
		}

		sp := strings.Split(cwd, string(filepath.Separator))
		return Resolve(strings.Join(sp[:len(sp)-1], sep), name)
	}

	goos := runtime.GOOS
	switch goos {
	case "linux", "darwin":
	default:
		return nil, errors.New(fmt.Sprintf("unsupported os %s", goos))
	}

	goarch := runtime.GOARCH
	switch goarch {
	case "amd64", "arm64":
	default:
		return nil, errors.New(fmt.Sprintf("unsupported arch %s", goarch))
	}

	files, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var so fs.FileInfo
	for _, file := range files {
		name := file.Name()
		if strings.Index(name, goos) != -1 && strings.Index(name, goarch) != -1 {
			so = file
			break
		}
	}
	if so == nil {
		return nil, errors.New(fmt.Sprintf("cannot find so with os: %s arch: %s", goos, goarch))
	}

	return plugin.Open(path.Join(p, so.Name()))
}
