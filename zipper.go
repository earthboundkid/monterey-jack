package main

import (
	"compress/gzip"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/carlmjohnson/monterey-jack/taskpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()
	return gzipAll(flag.Arg(0))
}

func gzipAll(root string) error {
	tp, _ := taskpool.New(context.Background(), runtime.NumCPU())

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		if info.IsDir() {
			if strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		for _, glob := range []string{"*.html", "*.htm", "*.css", "*.js", "*.svg"} {
			if matched, _ := filepath.Match(glob, base); matched {
				tp.Go(func() error { return gzipPath(path) })
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return tp.Wait()
}

func gzipPath(srcname string) error {
	destname := srcname + ".gz"

	// log.Printf("Gzipping %s to %s", srcname, destname)

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
