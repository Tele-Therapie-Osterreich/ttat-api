// Code generated by go-bindata. DO NOT EDIT.
// sources:
// migrations/0001-initial-schema.sql

package db


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataMigrations0001initialschemasql = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x58\x4f\x73\xea\x38\x12\xbf\xf3\x29\xba\xde\x25\x64\x16\x98\x3d\xed" +
	"\x61\x5e\x6d\x4d\x29\x58\x09\x9e\x67\xec\x3c\x5b\x4c\xc2\x5e\x28\x61\x37\x41\x3b\xc6\x66\x2d\x11\xc2\xb7\xdf\x6a" +
	"\xc9\x36\xc6\xf0\xfe\xec\x69\x7d\x09\x71\xab\xff\xfd\xba\x7f\xad\x86\xf1\x18\xfe\xb6\x53\x6f\x95\x34\x08\x8b\xfd" +
	"\x60\x30\x1e\xc3\xd7\x05\x4f\x84\x1f\x85\xbf\xc1\xcb\x8c\x89\xbb\x04\xc4\x8c\x43\xec\x3f\xcd\x04\xbc\xb0\x25\x88" +
	"\x08\x66\x2c\xf4\x02\x0e\x1e\x13\xec\x81\x25\x1c\x16\xcf\x1e\x13\x3c\x01\x16\x7a\x64\x60\xe6\x27\x22\x8a\x97\x30" +
	"\xe3\x31\xff\x1d\x92\x59\xb4\x08\xbc\xf6\xe5\x03\x07\x1e\xf8\x1e\xf7\x46\xe0\x4f\xf8\xc4\x1a\x27\x43\x20\x9c\x33" +
	"\x3a\xc5\xad\x19\x3f\x6c\x85\xd6\x8b\x9f\x00\x0b\x82\xe8\x85\x7b\x10\x85\x81\x0d\x64\xee\xc7\x71\x14\x03\x7f\x65" +
	"\x53\x11\x2c\x6d\xbc\x56\x67\x91\xf0\x98\x4c\x24\x6c\x99\x8c\x20\x8c\x04\x1d\xb6\x96\xeb\x38\xfc\x29\x0b\x60\xc8" +
	"\x42\x0f\xc2\x08\x82\x28\x7c\xe2\x31\x4c\xa3\x38\xe6\x53\x71\x6f\x3d\x4e\x06\xe3\x31\x59\x98\x2f\x41\xcc\xfc\xf0" +
	"\x8b\x1f\x3e\x01\x7b\x88\x16\x64\xdf\x4f\x28\x16\x8a\x17\xfc\x47\x58\x46\x0b\x78\x61\xa1\x75\xf1\xc0\x21\x11\xb1" +
	"\x3f\x15\xf5\x59\xb2\x44\x56\x9e\xe3\x48\xf0\x29\x81\x3a\xa2\xf8\x42\xab\xe4\x45\x36\x32\x17\xd6\x2f\x2c\x5c\xfe" +
	"\xd2\x0d\xce\x62\xf2\xc0\x97\x91\xc3\xb4\x2e\x05\x7f\x7d\x0e\xfc\xa9\x4f\xc9\xc6\xfc\xeb\xc2\x8f\xb9\x47\x7e\x9f" +
	"\xe3\xe8\x4f\xdf\xe3\x36\xf7\xc7\x45\x68\x3d\xb1\xc0\x17\x4b\x88\x1e\xed\xcb\xc4\x17\x9c\xcc\x88\x08\x7c\x91\x58" +
	"\x7c\x92\x89\x43\xd8\x4f\x60\x1a\x85\x82\xbf\x8a\x91\xcb\x69\xce\x59\xe8\xea\x40\x49\x87\x51\x27\x2a\x57\x96\xc7" +
	"\x28\x9e\x33\x72\x01\x4c\x50\x49\x26\xb6\x6d\x44\xb9\x1f\xe7\xf8\x8e\x39\x98\x2d\x56\x72\xaf\xb4\x81\x83\xc6\x0a" +
	"\xcc\x69\x8f\x7a\x32\x98\xc6\x9c\x09\x0e\x62\xf9\xcc\xcf\x27\x56\x24\x04\x96\x00\x0f\x17\x73\x18\xde\x95\xe6\x6e" +
	"\x04\x77\xfb\xed\x49\xab\x92\x3e\xe9\x3d\x62\xba\xbd\xbb\xff\x6c\x5d\xb0\xfd\xbe\x2a\xdf\x65\x0e\xda\x48\x73\xd0" +
	"\xb0\x29\xab\x8e\x33\x99\xa6\xe5\xa1\x30\x3d\x57\xb2\xd6\x59\x91\x4e\xd7\x55\x81\x47\xf2\xe0\xe4\x98\xd1\x67\xcc" +
	"\x94\xd1\xab\x3d\x16\x99\x2a\xde\xac\xfb\x83\xa6\xff\x30\x6b\x22\x88\x0a\x84\x0a\xd3\xb2\xca\x60\x8f\x55\x2f\xd3" +
	"\x09\x84\x78\xcc\x4f\x90\x56\x28\x0d\x66\x2e\xfb\x7d\x55\x6e\x54\x8e\x1a\x64\x85\x50\x94\x86\xac\xbc\x2b\xad\xd6" +
	"\x39\x82\x2a\x40\xa3\xac\xd2\x2d\x54\xa8\x0f\xb9\xd1\x70\x28\x8c\xb2\x00\x9e\xec\xf9\x26\xb8\x11\x1c\xb7\x2a\xdd" +
	"\xc2\x1b\x16\x58\xc9\x3c\x3f\x91\x95\x0a\xff\x73\x50\x15\x6a\x30\x5b\x69\x40\xc2\x4e\x15\x6a\x47\xe0\xa0\x81\x72" +
	"\x03\xaa\xd8\x94\xd5\x4e\x1a\x55\x16\xa0\x34\x6c\x54\x9e\x63\x46\x2e\x87\xa6\x84\x35\x3a\x0b\x64\xab\x51\x1c\x39" +
	"\x43\xc7\xf2\x90\x67\xf0\xef\x83\x36\xb0\x46\xc0\x9d\x54\xf9\x08\x0a\xb9\x43\x90\x45\xd6\xc9\x98\x2a\x77\xdf\xd0" +
	"\x44\x6c\x11\x64\x46\x86\xb4\xa9\xa4\x51\xef\x94\x9c\xc1\x6a\x23\x53\x04\xbd\x2d\x8f\x1a\x0a\x3c\x76\xc0\xd0\xf0" +
	"\xa9\xc6\x99\xd4\x9b\x22\x7d\x82\xb2\x48\x91\x9c\xf4\xb3\xc9\xa4\x91\x94\x86\x2a\xd2\xfc\x90\x61\xd6\x3a\x8e\xbc" +
	"\xe8\x37\x60\x9e\x07\x5f\xbe\xd8\x3e\x8e\xd9\x94\x38\xc5\xc4\x22\xb9\x94\x07\x4c\xfc\x1a\x44\x21\x04\xd1\xd4\x35" +
	"\x2f\x31\x6c\x04\x4f\x3c\x9a\x46\x1e\xf7\xa8\x97\x13\x2e\x16\xcf\x20\xfc\x39\xff\xbd\xed\x21\xf6\x10\x70\x5b\x49" +
	"\x0d\xc3\x01\x80\xca\xe0\xf2\x49\x78\xec\xb3\x00\x9e\x63\x7f\xce\xe2\x25\x7c\xe1\xcb\xd1\x00\x1c\x6c\xdd\x63\xc4" +
	"\x2f\x58\x84\xfe\xd7\x05\xb7\xac\x0f\x17\x41\x40\x07\x6d\xff\x77\x9f\x1e\x35\x9a\xb3\xe0\xf1\x47\xb6\x08\x04\x58" +
	"\x8a\x0c\xc0\x55\xa4\xef\x81\x04\xda\x54\x88\x66\x25\xb3\xac\x42\xad\x3b\x82\x54\x99\xd3\x4d\x8d\x7d\xa9\x4d\x5a" +
	"\x66\x78\x25\xb0\x7c\xaa\x4e\x37\x34\xb6\x65\x81\x37\x9d\x6f\xcb\xca\xac\xea\x32\x77\x05\x9b\x43\x9e\x77\xdf\x77" +
	"\xc2\xb5\x54\x3e\x3f\x3d\xc2\xb6\x79\x5b\xbe\x52\x50\x8e\x5d\x2b\x69\x1a\xdf\xfe\x9c\x27\x82\xcd\x9f\xc5\xbf\xae" +
	"\xe1\x2a\xca\xe3\xf0\x7e\x50\xb3\xf7\x65\x8b\x05\x48\xc7\x4b\xcb\x75\x42\x5b\xb5\x1c\x25\x02\x94\x1a\x6b\x51\xcd" +
	"\x57\x50\xbb\x1d\x66\x4a\x1a\x6c\x58\xb7\xc9\x31\x35\x8e\x4a\xd4\xa8\xfb\xc3\x3a\x57\x29\xbc\x2b\x3c\x52\xa3\x5e" +
	"\x58\x9c\x80\x5f\x68\x83\xd2\xf2\xc6\xd9\xb5\xed\x5e\x21\x6c\x31\xb7\x26\x24\x54\x65\xad\xa8\x34\x18\x49\x73\x41" +
	"\x6a\x90\xf0\x47\x12\x85\xe3\xb5\xd4\x98\xc1\x46\x61\x9e\xd5\xc3\x75\x2f\x4d\xba\x25\x23\x6b\x34\x47\xc4\xeb\x18" +
	"\x6a\x96\x12\xdf\xf2\x93\xf5\x89\x19\xbc\x63\xa5\x69\x10\xb8\x00\x49\xbd\x8d\xd0\x82\xd2\x86\x77\x31\x76\x60\x7d" +
	"\x02\x59\x74\x89\x5d\x56\x23\xe7\xb0\x89\x82\x26\x4b\x99\x67\x16\x0e\x53\x3a\x59\x5d\xe5\x1a\x9f\x9d\x54\x45\x4d" +
	"\x20\x9b\xdd\x04\x58\x71\x02\xfc\x30\x95\x3c\x23\xb2\xc6\x4d\xd9\x3a\x96\xb9\x8d\xe2\xca\x2e\xf9\x74\x56\x5b\xa4" +
	"\x74\x69\x87\x16\x99\x70\xe1\xa7\xb2\xa0\xb1\xd5\x66\x20\xf3\x1c\xa4\xb1\x63\x65\x72\xcd\xe9\x66\xd6\xaf\x9c\x72" +
	"\x8f\xe0\xb7\xa9\x6d\xf5\xdc\x29\x3f\x14\x9c\x36\x87\x1e\xaf\x21\xe6\x8f\x3c\xe6\xe1\x94\x27\x2e\xf1\xa1\xca\xee" +
	"\x81\x26\x0e\x0f\xb8\xe0\x30\x65\xc9\x94\x79\xdc\x12\xc9\xe6\x64\x1f\xaa\xf7\x83\x9d\x1d\xb6\x64\xb6\xbf\x7f\xbe" +
	"\xb5\xd9\xe1\x43\xe5\x4a\x56\xa7\x06\x18\x53\x56\xaa\x78\xbb\xb8\x83\x40\xed\xe4\x1b\x5d\xc5\xed\x5c\x9c\xce\x58" +
	"\xf8\xc4\xdd\x12\x20\x22\x5a\x0b\x80\x85\xc0\x5f\x05\x8f\x43\x16\x80\x3f\x67\x4f\x1c\xe6\x2c\x64\x4f\x7c\xce\x43" +
	"\x1a\x90\xf1\x9f\xfe\x94\x5f\x02\xe9\xac\xfe\x1f\xc0\xc3\x0f\x83\x85\x6d\x6b\x37\x5e\xbb\x73\xd5\xde\x17\xee\x79" +
	"\x58\x0a\xce\x1a\x9c\x68\x29\xe4\x01\x9f\x0a\x50\xd9\xa8\x7b\xbb\x8d\x6a\x69\xfd\x5c\x4e\xd1\x91\x1d\x9e\xa3\x76" +
	"\x52\x8e\x9a\xd1\x38\x72\xa3\xb0\xa7\xbc\xdf\x96\xa6\x1c\xf5\x86\xe1\x63\x1c\xcd\x6b\x1e\xbc\xd0\x66\x4c\x60\xf9" +
	"\x61\xa3\x38\xac\xa3\x6a\x00\x6a\x4f\xaf\xf4\x61\xbd\xd2\x7b\x4c\x95\xcc\x95\x51\xd8\x68\x5f\xbc\x3e\x91\xca\x3f" +
	"\xe1\xf7\x3a\xc9\xe4\xb0\x1e\x5f\xa8\xd0\xa2\x84\x32\xdd\xf6\xae\xef\x1e\x23\xae\x3c\xf5\x2f\xbc\xdb\x35\x3d\x47" +
	"\xd1\xbf\xbd\xfa\xf7\x54\x33\xf3\x73\xb9\xc6\xbc\xf7\x4e\xa5\x54\x47\xf7\x34\xad\xd1\x69\x05\xd7\x64\xbd\x5e\x48" +
	"\x78\xa7\xe0\x98\x23\x4d\x3a\x5b\xf1\x28\x0a\x38\x0b\x5b\xa6\x6c\x64\xae\x71\x34\x18\x40\xdd\x6c\xc3\x73\xc8\xae" +
	"\xf8\xf7\x9d\xf6\xf8\xa3\xa4\x11\x63\x49\x74\xb9\x5f\x52\x35\x7e\xd5\x5d\x68\x4f\xb0\x93\xc5\x69\x6c\xca\x31\xfd" +
	"\x85\x0a\x73\xbb\x6e\xdd\x1a\x34\x3f\xc4\xf6\x67\x69\xf3\x1d\x8c\x7e\x40\x97\xeb\x8e\xb9\x65\xa4\x1f\x68\xcf\x5e" +
	"\xcc\xdd\x17\x1c\x0b\x18\x21\xf6\x69\x2e\xdf\x54\x0a\xb9\x2a\xfe\xfa\x04\x79\xf9\x46\xe0\x95\x7f\x61\xa1\x7b\x28" +
	"\x58\xd1\xca\x89\x6c\xf2\xf6\x63\x9d\xcc\x74\xc6\xe2\xe1\x3f\xee\xdd\x3f\xdf\x59\xa5\x2c\xcd\x9b\xe7\xc6\x36\x95" +
	"\xcb\xe2\xed\x20\xdf\x2e\x36\x0f\xfc\xd8\xd3\x96\x7c\x35\x4f\x5b\x45\x4a\xa5\x0e\xd5\x0f\x3d\xfe\x7a\x11\xea\xca" +
	"\xa9\x67\x2b\x55\x64\xf8\x41\x40\x74\xa5\xc3\xb3\xf1\x86\x7a\xa8\xed\x44\xf2\x3d\xc7\x3a\x7b\xba\x8f\x85\x76\x87" +
	"\x7a\x38\xd8\xe4\xbe\x55\xf8\xff\xad\xdc\x6d\x75\x3c\xa9\xf2\x93\xdd\xb1\x94\x36\x2a\xd5\x20\xd7\xe5\xa1\xfe\x5e" +
	"\x66\x17\x29\x8a\x95\x36\x06\x4b\x9f\x6f\xf5\xee\x59\xfd\x46\xd7\xde\xee\xd8\x8c\x76\xb7\xf3\xe3\x91\xc9\xfe\xa4" +
	"\x6f\xe8\x99\x1e\xaa\x0a\x0b\xb3\x22\x1d\x3b\x32\xf0\xb8\x72\x83\xf2\xa2\xd3\x9b\xe3\x7f\xef\xd0\xbd\x39\x77\x75" +
	"\xa6\x21\x74\x0d\x80\xfb\x96\x75\x4e\xe4\xaa\x20\x24\xff\x41\xa2\x3f\x97\xa9\x4b\xf5\x87\x39\x5e\x92\xf1\x9c\xc0" +
	"\xad\x6b\xf0\x27\x29\x69\xcd\xda\x44\x50\xf7\xa7\x44\x07\xbb\x76\x0a\xc2\xd0\x06\xd3\x0b\xe5\xbe\xed\x9d\xf6\x77" +
	"\x21\xaf\x3c\x16\x83\x81\x17\x47\xcf\xdf\x82\xeb\x73\x57\xda\xeb\x99\xcf\x97\x9a\xae\xf3\x2f\x5e\x76\x19\x75\xc3" +
	"\x52\x2f\xfb\x4b\x7b\xdf\x13\x5e\x6f\x78\x57\xe2\xcb\x37\xee\x92\x69\x5e\x5d\xff\x7a\xd0\x95\x5c\xde\x74\x9f\x07" +
	"\xff\x0d\x00\x00\xff\xff\x9e\x90\x5b\x75\x49\x13\x00\x00")

func bindataMigrations0001initialschemasqlBytes() ([]byte, error) {
	return bindataRead(
		_bindataMigrations0001initialschemasql,
		"migrations/0001-initial-schema.sql",
	)
}



func bindataMigrations0001initialschemasql() (*asset, error) {
	bytes, err := bindataMigrations0001initialschemasqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "migrations/0001-initial-schema.sql",
		size: 4937,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1586207830, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}


//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"migrations/0001-initial-schema.sql": bindataMigrations0001initialschemasql,
}

//
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}


type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"migrations": {Func: nil, Children: map[string]*bintree{
		"0001-initial-schema.sql": {Func: bindataMigrations0001initialschemasql, Children: map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}