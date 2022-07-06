package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spasticus74/id3tagger/cli"
	"github.com/spasticus74/id3tagger/fileops"
	"github.com/spasticus74/id3tagger/parse"
)

func ProcessAlbum(dirPath, year, genre string) {
	// Pull out the name of all mp3 files in the dir
	mp3s, err := fileops.GetMP3sInDir(dirPath)
	if err != nil {
		fmt.Printf("A fatal error has occurred: %s", err)
	}

	// Stop if there's no files to process
	if len(mp3s) < 1 {
		fmt.Printf("No mp3s were found in the given path '%s'\n", dirPath)
		os.Exit(1)
	}

	// Parse out the artist and album names from the path
	splitpath := strings.Split(dirPath, string(filepath.Separator))
	albumName := splitpath[len(splitpath)-1]
	artistName := splitpath[len(splitpath)-2]
	artistName = cli.NamePrompt("Found Artist name:", artistName)
	albumName = cli.NamePrompt("Found Album name:", albumName)

	// If we haven't already entered Year and Genre prompt for them now
	if year == "" {
		year, err = cli.YearPrompt()
		if err != nil {
			fmt.Printf("Unable to parse an integer: %s\n", err)
		}
	}

	if genre == "" {
		genre = cli.GenrePrompt()
	}

	fmt.Printf("Using Artist: '%s', Album: '%s', Year: '%s', Genre: '%s'\n", artistName, albumName, year, genre)

	for _, v := range mp3s {
		trackNumber, trackName, err := parse.ParseFilename(v)
		if err != nil {
			fmt.Printf("Skipping '%s': %s\n", v, err)
		} else {
			fmt.Printf("Processing '%s':\t#:'%s', T:'%s'\n", v, trackNumber, trackName)
			err := Tag(dirPath+"/"+v, artistName, trackName, albumName, year, genre, trackNumber)
			if err != nil {
				fmt.Printf("An error occured in tagging '%s': '%s'. Continuing ...\n", v, err)
			}
		}

	}
}
