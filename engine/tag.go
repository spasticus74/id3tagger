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

	// Add track number
	ft := v2.V23FrameTypeMap["TRCK"]
	textFrame := v2.NewTextFrame(ft, track)
	mp3File.AddFrames(textFrame)

	err = mp3File.Close()
	if err != nil {
		return err
	}

	return nil
}
