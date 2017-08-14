package bindatafs

import (
	"errors"
	"io"
	"os"
	"syscall"
)

// ErrNotExist is the error returned when the file or directory does not exist
var ErrNotExist = os.ErrNotExist

// ErrOutOfBounds is the error returned by Read, ReadAt or Seek when reading/seeking past the end of the file
var ErrOutOfBounds = errors.New("Past boundaries of file")

// ErrInvalidWhence is returned by Seek when the given whence is not io.SeekStart, io.SeekCurrent or io.SeekEnd
var ErrInvalidWhence = errors.New("Specified whence is not supported")

type openFileFn func() []byte

type file struct {
	open   openFileFn
	data   []byte
	offset int
	size   int64
	info   os.FileInfo
}

func (f *file) prepare() {
	// prepare retrieves the content of the asset
	if nil == f.data {
		f.data = f.open()
		f.offset = 0
		f.size = int64(len(f.data))
	}
}

func (f *file) Close() error {
	// Close releases the content of the asset
	f.data = nil
	f.offset = 0
	f.size = 0
	return nil
}

func (f *file) Read(p []byte) (n int, err error) {
	f.prepare()
	if f.offset == len(f.data) {
		return 0, io.EOF
	}
	n = copy(p, f.data[f.offset:])
	f.offset += n
	return n, nil
}

func (f *file) ReadAt(p []byte, off int64) (n int, err error) {
	f.prepare()
	if off < 0 || off > f.size {
		return 0, ErrOutOfBounds
	}
	n = copy(p, f.data[int(off):])
	if n < len(p) {
		err = ErrOutOfBounds
	}
	return n, err
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	f.prepare()
	switch whence {
	case io.SeekStart:
		if offset > f.size {
			f.offset = int(f.size)
			return int64(f.offset), ErrOutOfBounds
		}
		if offset < 0 {
			f.offset = 0
			return 0, ErrOutOfBounds
		}
		f.offset = int(offset)
		return offset, nil
	case io.SeekCurrent:
		return f.Seek(int64(f.offset)+offset, io.SeekStart)
	case io.SeekEnd:
		return f.Seek(f.size+offset, io.SeekStart)
	default:
		return 0, ErrInvalidWhence
	}
}

func (f *file) Write(p []byte) (n int, err error) { return 0, syscall.EPERM }

func (f *file) WriteAt(p []byte, off int64) (n int, err error) { return 0, syscall.EPERM }

func (f *file) Name() string { return f.info.Name() }

func (f *file) Readdir(count int) ([]os.FileInfo, error) { return nil, ErrNotExist }

func (f *file) Readdirnames(n int) ([]string, error) { return nil, ErrNotExist }

func (f *file) Stat() (os.FileInfo, error) { return f.info, nil }

func (f *file) Sync() error { return nil }

func (f *file) Truncate(size int64) error { return syscall.EPERM }

func (f *file) WriteString(s string) (ret int, err error) { return 0, syscall.EPERM }
