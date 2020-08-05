package music

import (
	"fmt"
	"strconv"

	"github.com/HankiGreed/Innocent/pkg/config"
	"github.com/fhs/gompd/mpd"
)

type Music struct {
	Client *mpd.Client
}

func (m *Music) ConnectToClient() error {
	configs := config.ReadConfig()
	connectionString := configs.MPD.Address + ":" + strconv.Itoa(configs.MPD.Port)
	client, err := mpd.Dial("tcp", connectionString)
	m.Client = client
	return err
}

func (m *Music) ReturnPlaylists() []string {
	var playlistNames []string
	playlists, _ := m.Client.ListPlaylists()
	if len(playlists) == 0 {
		return append(playlistNames, "No Playlists Found")
	}
	for _, playlist := range playlists {
		playlistNames = append(playlistNames, playlist["playlist"])
	}
	return playlistNames
}

func (m *Music) ReturnSongsInPlaylist(name string) []string {
	var songNames []string
	songs, _ := m.Client.PlaylistContents(name)
	for i, song := range songs {
		songtext := fmt.Sprintf(" %3v %v ", i, song["Title"])
		songNames = append(songNames, songtext)
	}
	return songNames
}

func (m *Music) GetNowPlaying() string {
	status, err := m.Client.Status()
	if err != nil {
		return fmt.Sprintf("Error Occured : %s", err)
	}
	song, err := m.Client.CurrentSong()
	if err != nil {
		return fmt.Sprintf("Error Occured : %s", err)
	}
	songstring := " " + song["Title"] + " | " + song["Artist"] + " "
	if status["state"] == "pause" {
		songstring = " [Paused]" + songstring
	}

	return songstring
}

func (m *Music) ReturnStatusString() string {
	status, err := m.Client.Status()
	if err != nil {
		return fmt.Sprintf("Error Occured : %s", err)
	}

	statusString := ""
	if status["repeat"] == "1" {
		statusString += "Repeat : On, "
	} else {
		statusString += "Repeat : Off, "
	}

	if status["random"] == "1" {
		statusString += "Shuffle : On, "
	} else {
		statusString += "Shuffle : Off, "
	}
	statusString = statusString + "Volume : " + status["volume"]
	return statusString
}
