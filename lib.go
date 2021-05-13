package randomfiles

import (
	crand "crypto/rand"
	"io"
	"math/rand"
	"path"

	"github.com/colinmarc/hdfs"
)

type Options struct {
	Out    io.Writer // output progress
	Source io.Reader // randomness source

	FileSize int // the size per file.

	Depth int // how deep the hierarchy goes
	Files int // how many files per dir
	Width int // how many dirs per dir

	RandomFanout bool // randomize fanout numbers
}

var FilenameSize = 16
var alphabet = []rune("abcdefghijklmnopqrstuvwxyz01234567890-_")

func WriteRandomFiles(root string, depth int, opts *Options, client *hdfs.Client) error {
	numFiles := opts.Files

	for i := 0; i < numFiles; i++ {
		if e := WriteRandomFile(root, opts, client); e != nil {
			return e
		}
	}

	if depth+1 <= opts.Depth {
		numDirs := opts.Depth
		for i := 0; i < numDirs; i++ {
			if e := WriteRandomDir(root, depth+1, opts, client); e != nil {
				return e
			}
		}
	}

	return nil
}

func WriteRandomFile(root string, opts *Options, client *hdfs.Client) error {
	filesize := int64(opts.FileSize)

	n := rand.Intn(FilenameSize-4) + 4
	name := RandomFilename(n)
	filepath := path.Join(root, name)
	f, e := client.Create(filepath)
	if e != nil {
		return e
	}

	if _, e := io.CopyN(f, crand.Reader, filesize); e != nil {
		return e
	}
	return f.Close()
}

func RandomFilename(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func WriteRandomDir(root string, depth int, opts *Options, client *hdfs.Client) error {
	if depth > opts.Depth {
		return nil
	}
	n := rand.Intn(FilenameSize-4) + 4
	name := RandomFilename(n)
	root = path.Join(root, name)
	if e := client.MkdirAll(root, 0755); e != nil {
		return e
	}

	return WriteRandomFiles(root, depth, opts, client)
}
