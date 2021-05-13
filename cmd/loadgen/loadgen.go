package main

import (
	"flag"
	"fmt"
	"os"

	randomfiles "github.com/Galzzly/loadgenerator"
	"github.com/colinmarc/hdfs"
)

var usage = `Usage: %s [OPTIONS] <path>
Write a random filesystem heirarchy to each <path> in HDFS.

Options:
`

var opts randomfiles.Options
var paths []string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&opts.FileSize, "filesize", 1000000, "filesize - max file size")
	flag.IntVar(&opts.Depth, "depth", 3, "depth - how deep to you want the directory tree")
	flag.IntVar(&opts.Width, "width", 2, "width - number of subdirectories per directory")
	flag.IntVar(&opts.Files, "files", 15, "files - total number of files")
}

func parseArgs() error {
	flag.Parse()

	paths = flag.Args()
	if len(paths) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	return nil
}

func run() {
	/*
		Launch the HDFS client connection.
		Not passing a connection string will attempt to find the details from
		config files
	*/
	client, e := hdfs.Client("") // Should get the connection from the config files
	if e != nil {
		return e
	}

	for _, root := range paths {
		if e := client.MkdirAll(root, 0755); e != nil {
			return e
		}

		e := randomfiles.WriteRandomFiles(root, 1, &opts, client)
		if e != nil {
			return e
		}
	}
	return nil
}

func main() {
	// Parse the arguments
	if e := parseArgs(); e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		os.Exit(1)
	}

	// Let's start then
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		os.Exit(1)
	}
}
