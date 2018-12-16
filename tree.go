package bindatafs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type tree struct {
	root *dir
}

func sanitize(name string) string {
	if strings.Contains(name, "/") {
		return name
	}
	return strings.Replace(name, "\\", "/", -1)
}

func splitList(name string) []string {
	name = sanitize(name)
	names := strings.Split(name, "/")
	var i = len(names)
	// Right trim empty path parts
	for ; i > 0; i-- {
		if "" != names[i-1] {
			break
		}
	}
	return names[:i]
}

func split(name string) (string, string) {
	name = sanitize(name)
	last := strings.LastIndex(name, "/")
	if -1 == last {
		return "", name
	}
	return name[:last], name[last+1:]
}

func (t *tree) dir(name string, create bool) *dir {
	dir := t.root
	for _, d := range splitList(name) {
		if "" == d {
			continue
		}
		dir = dir.DirByName(d, create)
		if nil == dir {
			return nil
		}
	}
	return dir
}

func (t *tree) grow(fs *Fs) error {
	t.root = &dir{
		info: &fileInfo{"", 0, os.ModeDir, time.Unix(0, 0), true},
	}
	names := fs.Names()
	for _, name := range names {
		dirName, fileName := split(name)
		dir := t.dir(dirName, true)
		info, err := fs.Info(name)
		if nil != err {
			return err
		}
		// Makes assumption that asset can be retrieved as well when the info can be retrieved
		// Create new info because go-bindata returns full path instead of file name for it's Name() function
		var nameCapture = name
		dir.files = append(dir.files, &file{
			open: func() []byte { return fs.Asset(nameCapture) },
			info: &fileInfo{fileName, info.Size(), info.Mode(), info.ModTime(), false},
		})
	}
	return nil
}

func (t *tree) Open(fs *Fs, name string) (afero.File, error) {
	if nil == t.root {
		err := t.grow(fs)
		if nil != err {
			return nil, fmt.Errorf("Unable to initialize file tree; %v", err)
		}
	}
	dirName, fileName := split(name)
	dir := t.dir(dirName, false)
	if nil == dir {
		return nil, ErrNotExist
	}
	if "" == fileName {
		return dir, nil
	}
	file := dir.FileByName(fileName)
	if nil == file {
		return nil, ErrNotExist
	}
	return file, nil
}
