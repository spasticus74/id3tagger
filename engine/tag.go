package engine

import (
	"fmt"

	"github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
)

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

	// Remove any already existing TRCK tags
	trcks := mp3File.Frames("TRCK")
	if len(trcks) > 0 {
		mp3File.DeleteFrames("TRCK")
	}
	// Add one back in
	ft := v2.V23FrameTypeMap["TRCK"]
	textFrame := v2.NewTextFrame(ft, track)
	mp3File.AddFrames(textFrame)

	err = mp3File.Close()
	if err != nil {
		return err
	}

	return nil
}

func Inspect(mp3Filepath string) error {
	mp3File, err := id3.Open(mp3Filepath)

	if err != nil {
		return err
	}
	defer mp3File.Close()

	f := mp3File.AllFrames()
	fmt.Printf("Found %d frames in %s\n", len(f), mp3Filepath)

	for c, v := range f {
		fmt.Printf("%d: %s: %s\n", c+1, v.Id(), v.String())
	}
	fmt.Println()

	return nil
}
