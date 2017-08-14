package bindatafs

import (
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

type dir struct {
	info  os.FileInfo
	dirs  []*dir
	files []*file
}

// FileByName returns either a file of a sub directory by it's name.
func (d *dir) FileByName(name string) afero.File {
	for _, f := range d.files {
		if name == f.Name() {
			return f
		}
	}
	for _, sub := range d.dirs {
		if name == sub.Name() {
			return sub
		}
	}
	return nil
}

// DirByName returns a sub directory (non-recursive) by it's name. If create is true then the directory is created when it doesn't exist yet
func (d *dir) DirByName(name string, create bool) *dir {
	for _, sub := range d.dirs {
		if name == sub.Name() {
			return sub
		}
	}
	if create {
		// go-bindata doesn't store mod time for directories
		dir := &dir{info: &fileInfo{name, 0, os.ModeDir, time.Unix(0, 0), true}}
		d.dirs = append(d.dirs, dir)
		return dir
	}
	return nil
}

func (d *dir) Close() error { return nil }

func (d *dir) Read(p []byte) (n int, err error) { return 0, syscall.EPERM }

func (d *dir) ReadAt(p []byte, off int64) (n int, err error) { return 0, syscall.EPERM }

func (d *dir) Seek(offset int64, whence int) (int64, error) { return 0, syscall.EPERM }

func (d *dir) Write(p []byte) (n int, err error) { return 0, syscall.EPERM }

func (d *dir) WriteAt(p []byte, off int64) (n int, err error) { return 0, syscall.EPERM }

func (d *dir) Name() string { return d.info.Name() }

func (d *dir) Readdir(n int) ([]os.FileInfo, error) {
	dirsN := len(d.dirs)
	filesN := len(d.files)
	max := dirsN + filesN
	if n < 0 || n > max {
		n = max
	}
	files := make([]os.FileInfo, n)
	for i := 0; i < dirsN; i++ {
		files[i] = d.dirs[i].info
	}
	for i := 0; i < filesN; i++ {
		files[i] = d.files[i-dirsN].info
	}
	return files, nil
}

func (d *dir) Readdirnames(n int) ([]string, error) {
	dirsN := len(d.dirs)
	filesN := len(d.files)
	max := dirsN + filesN
	if n < 0 || n > max {
		n = max
	}
	files := make([]string, n)
	for i := 0; i < dirsN; i++ {
		files[i] = d.dirs[i].info.Name()
	}
	for i := 0; i < filesN; i++ {
		files[i+dirsN] = d.files[i].info.Name()
	}
	return files, nil
}

func (d *dir) Stat() (os.FileInfo, error) { return d.info, nil }

func (d *dir) Sync() error { return nil }

func (d *dir) Truncate(size int64) error { return syscall.EPERM }

func (d *dir) WriteString(s string) (ret int, err error) { return 0, syscall.EPERM }
