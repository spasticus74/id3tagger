package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	id3 "github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v2"
)

var dirPath string

func main() {
	dirPtr := flag.String("p", ".", "path to the instructions file")

	flag.Parse()

	// Clean up the path we've receieved because we need to be able to
	// do some reliable parsing later if it's valid
	dirPath, _ = filepath.Abs(*dirPtr)
	dirPath = filepath.Clean(dirPath)
	fmt.Printf("Processing path %s ...\n", dirPath)

	// First validate the dir path
	err := CheckPath(dirPath)
	if err != nil {
		fmt.Printf("A fatal error has occurred: %s\n", err)
	}

	// Pull out the name of all mp3 files in the dir
	mp3s, e := GetMP3sInDir(*dirPtr)
	if e != nil {
		fmt.Printf("A fatal error has occurred: %s", e)
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
	artistName = NamePrompt("Found Artist name:", artistName)
	albumName = NamePrompt("Found Album name:", albumName)
	albumYear, err := YearPrompt()
	if err != nil {
		fmt.Printf("Unable to parse an integer: %s\n", err)
	}
	albumGenre := GenrePrompt()

	fmt.Printf("Using Artist: '%s', Album: '%s', Year: '%s', Genre: '%s'\n", artistName, albumName, albumYear, albumGenre)

	for _, v := range mp3s {
		trackNumber, trackName, err := ParseFilename(v)
		if err != nil {
			fmt.Printf("Skipping '%s': %s\n", v, err)
		} else {
			fmt.Printf("Processing '%s':\t#:'%s', T:'%s'\n", v, trackNumber, trackName)
			err := Tag(dirPath + "/" + v, artistName, trackName, albumName, albumYear, albumGenre, trackNumber)
			if err != nil {
				fmt.Printf("An error occured in tagging '%s': '%s'. Continuing ...\n", v, err)
			}
		}

	}
}

func Tag(mp3Filepath, artist, title, album, year, genre, track string) error {
	mp3File, err := id3.Open(mp3Filepath)

	if err != nil {
		return err
	}
	defer mp3File.Close()

	tag := mp3File.Tagger
	if tag == nil {
		fmt.Println("error: no tag added to file")
	}
	mp3File.SetArtist(artist)
	mp3File.SetTitle(title)
	mp3File.SetAlbum(album)
	mp3File.SetYear(year)
	mp3File.SetGenre(genre)

	// track number
	ft := v2.V23FrameTypeMap["TRCK"]
	textFrame := v2.NewTextFrame(ft, track)
	mp3File.AddFrames(textFrame)

	err = mp3File.Close()
	if err != nil {
		return err
	}

	return nil
}

func GenrePrompt() string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Enter a genre, or press enter to ignore: ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	return ret
}

func YearPrompt() (string, error) {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Enter a year, or press enter to ignore: ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	if strings.ToUpper(ret) == "" {
		return "-1", nil
	} else {
		_, e := strconv.Atoi(ret)
		if e != nil {
			return "-1", e
		}
		return ret, nil
	}
}

func NamePrompt(label, value string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s %s. Enter new value or press enter to accept: ", label, value)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	if strings.ToUpper(ret) == "" {
		return value
	} else {
		return ret
	}
}

func CheckPath(dirPath string) error {
	// first, check that what we've been given actually exists
	theDir, err := os.Open(dirPath)

	if err != nil {
		theDir.Close()
		return err
	} else {
		fileInfo, err := theDir.Stat()
		if err != nil {
			theDir.Close()
			return err
		}

		// check that the path given points to a directory and not a file or some other thing
		if !fileInfo.IsDir() {
			theDir.Close()
			return errors.New("the given path exists but is not a directory")
		}
	}

	// all clear
	theDir.Close()
	return nil
}

func GetMP3sInDir(dirPath string) ([]string, error) {
	// open the dir
	theDir, _ := os.Open(dirPath)
	defer theDir.Close()

	// read the dir
	files, err := theDir.Readdir(0)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	// gather the ones that are mp3
	var ret []string
	for _, v := range files {
		if filepath.Ext(v.Name()) == ".mp3" {
			ret = append(ret, v.Name())
		}
	}

	return ret, nil
}

func ParseFilename(filename string) (string, string, error) {
	re := regexp.MustCompile(`(\d*) (.*).mp3`)
	matches := re.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return "-1", "", errors.New("unable to parse filename")
	}

	n, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return "-1", "", errors.New("unable to parse filename")
	}

	return strconv.Itoa(n), matches[0][2], nil
}
