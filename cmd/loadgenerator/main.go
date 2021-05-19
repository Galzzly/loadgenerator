package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	lg "github.com/Galzzly/loadgenerator"
	"github.com/Galzzly/loadgenerator/randomfiles"
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
	/*
		Keeping these areound just in case
		flag.StringVar(&opts.Keytab, "keytab", "", "path to the keytab file")
		flag.StringVar(&opts.Principal, "principal", "", "the principal to use with the keytab")
		flag.StringVar(&opts.Realm, "realm", "", "the realm to use with the keytab")
	*/
}

func parseArgs() error {
	flag.Parse()

	paths = flag.Args()
	if len(paths) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	/*
		Want to check for the presence of either of the os env variables:
			HADOOP_HOME
			HADOOP_CONF_DIR
		If neither are set, then set HADOOP_HOME to the standard location
		of `/etc/hadoop`

		This will save the user having to launc using:
			HADOOP_HOME=/etc/hadoop ./loadgenerator {path}
	*/
	if os.Getenv("HADOOP_HOME") == "" && os.Getenv("HADOOP_CONF_DIR") == "" {
		if e := os.Setenv("HADOOP_HOME", "/etc/hadoop"); e != nil {
			return e
		}
	}

	return nil
}

func run() error {
	client, e := lg.ConnectToNamenode()
	if e != nil {
		return e
	}

	for _, root := range paths {
		fmt.Printf("Generating tree for %s ...", root)
		if e := client.MkdirAll(root, 0755); e != nil {
			return e
		}

		e := randomfiles.WriteRandomFiles(root, 1, &opts, client)
		if e != nil {
			return e
		}
		fmt.Println("Done")
	}

	return nil
}

func main() {
	start := time.Now()
	// Parse the arguments
	if e := parseArgs(); e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		os.Exit(1)
	}

	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		os.Exit(1)
	}

	fmt.Println("Time Taken:", time.Since(start))
}
