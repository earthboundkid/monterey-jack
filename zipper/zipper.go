package zipper

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/carlmjohnson/monterey-jack/taskpool"
)

// All calls FromPath for all files in root matching a glob.
func All(root string, globs ...string) error {
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

		for _, glob := range globs {
			if matched, _ := filepath.Match(glob, base); matched {
				tp.Go(func() error { return FromPath(path) })
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

// FromPath gzips a file from the given source pathname to that path plus ".gz".
func FromPath(srcname string) error {
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
