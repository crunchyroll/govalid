// chris 071415

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var args struct {
	// Program name variables.  Set by init.
	prog      string
	progUpper string

	// Source and destination file names.  Can be "-" for stdin/out.
	srcname string
	dstname string
}

func usage() {
	log.Printf("usage: %s [-h] [-o file.go] [file.v]", args.prog)
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	args.prog      = path.Base(os.Args[0])
	args.progUpper = strings.ToUpper(args.prog)
	log.SetFlags(0)
	flag.Usage = usage
	dstname := flag.String("o", "", "output file name")
	flag.Parse()

	if len(flag.Args()) == 0 {
		args.srcname = "-"
	} else if len(flag.Args()) == 1 {
		args.srcname = flag.Args()[0]
	} else {
		usage()
	}

	args.dstname = *dstname
}

func main() {
	var (
		src io.Reader
		dst io.Writer
	)

	if args.srcname == "-" {
		src = os.Stdin
	} else {
		file, err := os.Open(args.srcname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		src = file
	}

	if args.dstname == "" {
		// Need to determine where output will go based on the
		// source.
		if args.srcname == "-" {
			args.dstname = "-"
		} else {
			// Chop off srcname extension and replace it
			// with ".go".
			srcext := path.Ext(args.srcname)
			args.dstname = fmt.Sprintf("%s.go", args.srcname[:len(args.srcname)-len(srcext)])
		}
	}

	if args.dstname == "-" {
		dst = os.Stdout
	} else {
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		file, err := os.OpenFile(args.dstname, flag, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		dst = file
	}

	if err := process(dst, args.srcname, src); err != nil {
		log.Fatal(err)
	}
}
