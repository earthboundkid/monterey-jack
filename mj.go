package main

import (
	"flag"
	"log"

	"github.com/carlmjohnson/monterey-jack/zipper"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()
	return zipper.All(flag.Arg(0),
		".html", ".htm", ".css", ".js", ".svg", ".xml")
}
