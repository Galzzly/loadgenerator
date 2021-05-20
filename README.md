# loadgenerator ![issues](https://img.shields.io/github/issues/Galzzly/loadgenerator?style=plastic) [![extract GoDoc](https://img.shields.io/badge/reference-godoc-blue.svg?logo=go&style=plastic)](https://pkg.go.dev/github.com/Galzzly/loadgenerator)

Introducing **loadgenerator v1** - a small utility to generate a directory directly into HDFS (rather than generating one on a local filesystem and uploading).

## Features

Package loadgenerator attempts to make it quick and easy, without the need to launch a java process.

The `loadgenerator` command can be ran simply as `./loadgenerator /path` to generate a tree under `/path` using the defaults.

The utility makes use of one of the `HADOOP_HOME` or `HADOOP_CONF_DIR` environment variable. If this is not set, it will use a default of `HADOOP_HOME=/etc/hadoop`. Should a more specific location be required, then either export the appropriate environment variable to necessary location, or run the command as `HADOOP_CONF_DIR=/path/to/conf ./loadgenerator /path`.

### Usage
```
Usage: ./loadgenerator [OPTIONS] <path>
Write a random filesystem heirarchy to each <path> in HDFS.

Options:
  -depth int
    	depth - how deep to you want the directory tree (default 3)
  -files int
    	files - total number of files (default 15)
  -filesize int
    	filesize - max file size (default 1000000)
  -width int
    	width - number of subdirectories per directory (default 2)
```

## GoDoc

See <https://pkg.go.dev/github.com/Galzzly/loadgenerator>

## Install

### GO

To install the binary directly into your \$GOPATH/bin:

```
go get github.com/galzzly/loadgenerator/cmd/loadgenerator
```
### Binary

The latest compiled binary can be found here <https://github.com/Galzzly/loadgenerator/releases>

## Issues

Currently, loadgenerator will not work on Kerberised environments. The code is currently being worked on. One of the issues is when weak ciphers are being used. The

## Thanks

Thanks to <a href="https://github.com/jbenet/go-random-files">jbenet/go-random-files</a> for the idea and reference point for the methodology to build out the structure.

Thanks to <a href="https://github.com/colinmarc/hdfs">colinmarc/hdfs</a> for building a HDFS utility to connect directly into a Hadoop cluster. 