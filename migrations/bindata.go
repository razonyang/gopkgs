package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __20200712083244_init_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\x48\xc9\xcf\x4d\xcc\xcc\x2b\x4e\xb0\xe6\x42\x16\x2d\x48\x4c\xce\x4e\x4c\x4f\x45\x17\x4e\x4c\x2e\xc9\xcc\xc7\x50\x9c\x9c\x98\x93\x9a\x97\x92\x58\x04\x12\x07\x04\x00\x00\xff\xff\x78\xda\x24\x3f\x5b\x00\x00\x00")

func _20200712083244_init_down_sql() ([]byte, error) {
	return bindata_read(
		__20200712083244_init_down_sql,
		"20200712083244_init.down.sql",
	)
}

var __20200712083244_init_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x93\xb1\x6e\xea\x30\x14\x86\xf7\x3c\xc5\x19\x1d\x89\xe1\x5e\xee\x6d\x17\x26\x03\x6e\x1b\x35\x98\x36\x72\x2a\x31\x39\x56\xec\x82\x05\x24\x28\x84\x8a\xc7\xaf\x82\xdd\xc4\x90\x10\x24\x86\xae\xf1\x77\x8e\x7f\x7f\x3f\x4c\x22\x82\x19\x01\x86\xc7\x21\x81\x44\xe6\x5b\xa1\xb3\x7d\x02\xc8\x03\x00\x48\xb4\x4c\x60\x1c\x3c\x07\x94\x01\x9d\x33\xa0\x71\x18\x02\x8e\xd9\x9c\x07\x74\x12\x91\x19\xa1\x6c\x60\xc0\xc3\x5e\x15\xbc\xa2\x3f\x70\x34\x79\xc1\x11\x7a\xfc\xef\xd7\x23\x96\xc9\xc4\x56\xf5\x02\x5f\xaa\xd0\x9f\x5a\xc9\x04\x58\x40\x17\x67\x97\x4e\xc9\x13\x8e\x43\x06\x7f\x2c\x9a\xae\xc4\x66\xa3\xb2\xa5\xe2\xe5\xb1\x6c\x96\xfe\x1b\xb6\x96\xa6\x85\x12\xa5\x92\x5c\x94\x09\x4c\x31\x23\x2c\x98\x91\x4b\xe6\xb0\x93\xb7\x98\xb7\x28\x98\xe1\x68\x01\xaf\x64\x81\x2a\x2d\xbe\xf9\x1c\xd3\xe0\x3d\x26\xd5\xd7\x4a\xd6\x91\x5b\x81\xdc\xfa\xe0\xe6\xcd\xa8\xf6\x33\xb0\x1a\xec\x78\x6b\xae\x51\x80\x1a\x1d\xbe\xe7\x8f\x3c\xef\xbc\xa9\x9d\x48\xd7\x62\xa9\xee\xa8\xca\x5c\xc5\x3b\x78\x0b\xec\x44\xb9\xea\xef\x29\xdd\x3b\xe7\xad\xe3\x22\xcf\x9d\x4a\x86\x0f\x6d\x42\xe6\xee\x86\x2e\xe2\x77\x5b\xfb\x91\xc9\x6b\x37\xdc\x48\x40\x8e\xad\x81\x35\xd3\x54\x77\x3e\x7b\xb2\x82\x4e\x72\xba\x0a\x13\x69\xa9\xf3\x8b\xbf\x56\xcf\xcf\xd6\xae\xed\xa9\x69\xad\x33\x67\xc5\xdf\xbb\x1c\x3a\x7e\xc0\x15\x54\x9b\xb1\xa9\xb9\x1b\x07\xb9\xe1\xae\x0d\x98\x74\xc8\xa4\xbc\x06\xb9\x09\x91\x9b\xb7\x4b\x60\x2a\x36\x2a\x93\xa2\xa8\x15\x6a\x79\x7a\xd6\xad\xca\xc1\xf3\x47\xdf\x01\x00\x00\xff\xff\x69\xba\x60\xca\xe6\x04\x00\x00")

func _20200712083244_init_up_sql() ([]byte, error) {
	return bindata_read(
		__20200712083244_init_up_sql,
		"20200712083244_init.up.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"20200712083244_init.down.sql": _20200712083244_init_down_sql,
	"20200712083244_init.up.sql": _20200712083244_init_up_sql,
}
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
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"20200712083244_init.down.sql": &_bintree_t{_20200712083244_init_down_sql, map[string]*_bintree_t{
	}},
	"20200712083244_init.up.sql": &_bintree_t{_20200712083244_init_up_sql, map[string]*_bintree_t{
	}},
}}
