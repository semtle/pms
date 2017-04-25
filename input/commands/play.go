package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"

	"github.com/fhs/gompd/mpd"
)

// Play plays songs in the MPD playlist.
type Play struct {
	songlistWidget *widgets.SonglistWidget
	mpdClient      func() *mpd.Client
	song           *song.Song
	id             int
	pos            int
}

func NewPlay(songlistWidget *widgets.SonglistWidget, mpdClient func() *mpd.Client) *Play {
	return &Play{songlistWidget: songlistWidget, mpdClient: mpdClient}
}

func (cmd *Play) Reset() {
	cmd.song = nil
	cmd.pos = -1
}

func (cmd *Play) Execute(t lexer.Token) error {
	var err error

	s := t.String()

	switch t.Class {
	case lexer.TokenIdentifier:
		switch s {
		case "cursor":
			cmd.song = cmd.songlistWidget.CursorSong()
			if cmd.song == nil {
				return fmt.Errorf("Cannot play: no song under cursor")
			}
		default:
			return nil
		}

	case lexer.TokenEnd:
		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		if cmd.song == nil {
			err = client.Play(-1)
			return err
		}

		// Add song to queue only if we are not operating on the queue
		id := cmd.song.ID
		list := cmd.songlistWidget.Songlist()

		switch list.(type) {
		case *songlist.Queue:
		default:
			id, err = client.AddID(cmd.song.TagString("file"), -1)
			if err != nil {
				return err
			}
		}

		err = client.PlayID(id)
		return err

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return nil
}
