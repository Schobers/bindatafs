package bindatafs

import (
	"os"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (i *fileInfo) Name() string { return i.name }

func (i *fileInfo) Size() int64 { return i.size }

func (i *fileInfo) Mode() os.FileMode { return i.mode }

func (i *fileInfo) ModTime() time.Time {
	return i.modTime
}

func (i *fileInfo) IsDir() bool { return i.isDir }

func (i *fileInfo) Sys() interface{} { return nil }
