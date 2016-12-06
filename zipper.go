package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
)

func run() error {
	flag.Parse()
	srcname := flag.Arg(0)
	destname := srcname + ".gz"

	log.Printf("Gzipping %s to %s", srcname, destname)

	src, err := os.Open(srcname)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destname)
	if err != nil {
		return err
	}
	defer dest.Close()

	w, _ := gzip.NewWriterLevel(dest, gzip.BestCompression)
	_, err = io.Copy(w, src)
	if err != nil {
		return err
	}

	return w.Close()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
