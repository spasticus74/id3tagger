package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spasticus74/id3tagger/engine"
	"github.com/spasticus74/id3tagger/fileops"
)

func main() {
	dirPtr := flag.String("p", ".", "path to the album directory")
	yearPtr := flag.String("y", "", "year of release")
	genrePtr := flag.String("g", "", "album genre")
	inspPtr := flag.Bool("i", false, "inspect the existing tags")

	flag.Parse()

	albumYear := *yearPtr
	albumGenre := *genrePtr

	// Clean up the path we've receieved because we need to be able to
	// do some reliable parsing later if it's valid
	dirPath, _ := filepath.Abs(*dirPtr)
	dirPath = filepath.Clean(dirPath)
	fmt.Printf("Processing path %s ...\n", dirPath)

	// First validate the dir path
	err := fileops.CheckPath(dirPath)
	if err != nil {
		fmt.Printf("A fatal error has occurred: %s\n", err)
		os.Exit(1)
	}

	if *inspPtr {
		// Report the number of frames in each file
		engine.InspectAlbum(dirPath)
	} else {
		// Process the dir
		engine.ProcessAlbum(dirPath, albumYear, albumGenre)
	}
}
