# bindatafs

Bindatafs combines the [Afero file system abstraction][0] with [file embedding of go-bindata][1].

## How to use
### Getting it
```
go get github.com/Schobers/bindatafs
```

### Using it
1. Install go-bindata
2. Generate your bindata using go-bindata
3. Create an bindatafs file system instance:
```
var fs afero.Fs = bindatafs.NewFs(MustAsset, AssetInfo, AssetNames)
```
or
```
var fs afero.Fs = &bindatafs.Fs{Asset: MustAsset, Info: AssetInfo, Names: AssetNames}
```

### Forks of go-bindata
Bindatafs works with any fork of go-bindata that provides the same API as the [original][1], namely the methods MustAsset, AssetInfo and AssetNames. Tested with:
- [https://github.com/shuLhan/go-bindata][2]


[0]: https://github.com/spf13/afero
[1]: https://github.com/jteeuwen/go-bindata
[2]: https://github.com/shuLhan/go-bindata