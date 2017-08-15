package bindatafs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fsCSSContent = []byte("body {\r\n  color: #0c8add81;\r\n}")
var binJSContent = []byte("function bin() {\r\n  return fs();\r\n}")
var fsJSContent = []byte("function fs() {\r\n  return \"afero\";\r\n}")

func createFs() *Fs {
	return &Fs{Asset: MustAsset, Info: AssetInfo, Names: AssetNames}
}

func read(fs *Fs, tree *tree, name string) ([]byte, error) {
	f, err := tree.Open(fs, name)
	if nil != err {
		return nil, err
	}
	info, err := f.Stat()
	if nil != err {
		return nil, err
	}
	size := info.Size()

	var content = make([]byte, size)
	_, err = f.Read(content)
	if nil != err {
		return nil, err
	}
	return content, nil
}

func TestTreeGrow(t *testing.T) {
	var tree tree
	var fs = createFs()
	f, err := tree.Open(fs, "")
	require.Nil(t, err)
	assert.NotNil(t, f)
}

func TestOpenNonExistingFile(t *testing.T) {
	var tree tree
	var fs = createFs()
	_, err := tree.Open(fs, "tests/css/nofs.css")
	require.NotNil(t, err)
}

func TestOpenFile(t *testing.T) {
	var tree tree
	var fs = createFs()
	f, err := tree.Open(fs, "tests/css/fs.css")

	require.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "fs.css", f.Name())
}

func TestStatOnOpenFile(t *testing.T) {
	var tree tree
	var fs = createFs()
	f, _ := tree.Open(fs, "tests/css/fs.css")

	info, err := f.Stat()
	require.Nil(t, err)
	assert.Equal(t, int64(len(fsCSSContent)), info.Size())
	assert.False(t, info.IsDir())
}

func TestReadOnOpenFile(t *testing.T) {
	var fs = createFs()
	var tree tree

	f, _ := tree.Open(fs, "tests/css/fs.css")
	info, _ := f.Stat()
	size := info.Size()

	var content = make([]byte, size)
	n, err := f.Read(content)
	require.Nil(t, err)
	assert.Equal(t, int(size), n)
	assert.Equal(t, fsCSSContent, content)
}

func TestReadOnFile(t *testing.T) {

	var fs = createFs()
	var tree tree

	assertContent := func(t *testing.T, err error, expected []byte, actual []byte) {
		require.Nil(t, err)
		assert.Equal(t, expected, actual)
	}

	content, err := read(fs, &tree, "tests/css/fs.css")
	assertContent(t, err, fsCSSContent, content)
	content, err = read(fs, &tree, "tests/js/bin.js")
	assertContent(t, err, binJSContent, content)
	content, err = read(fs, &tree, "tests/js/fs.js")
	assertContent(t, err, fsJSContent, content)
}
