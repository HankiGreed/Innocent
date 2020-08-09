package music

import (
	"fmt"
	"log"
	"strconv"

	"github.com/HankiGreed/Innocent/pkg/config"
	"github.com/fhs/gompd/mpd"
)

type Music struct {
	Client          *mpd.Client
	Repeat, Shuffle bool
}

func (m *Music) ConnectToClient() error {
	configs := config.ReadConfig()
	connectionString := configs.MPD.Address + ":" + strconv.Itoa(configs.MPD.Port)
	client, err := mpd.Dial("tcp", connectionString)
	m.Client = client
	status, err := m.Client.Status()
	if err != nil {
		log.Fatalln(err)
	}
	repeat, err := strconv.ParseBool(status["repeat"])
	m.Repeat = repeat
	shuffle, err := strconv.ParseBool(status["random"])
	m.Shuffle = shuffle
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

func (m *Music) GetNowPlaying() (string, int, string) {
	status, err := m.Client.Status()
	if err != nil {
		log.Fatalln(err)
	}
	song, err := m.Client.CurrentSong()
	if err != nil {
		log.Fatalln(err)
	}
	if status["state"] == "stop" {
		return " Nothing Playing ", 0, "No Progress"
	}
	songstring := " " + song["Title"] + " | " + song["Artist"] + " "
	if status["state"] == "pause" {
		songstring = " [Paused]" + songstring
	}
	elapsed, err := strconv.ParseFloat(status["elapsed"], 64)
	if err != nil {
		log.Fatalln(err)
	}
	total, err := strconv.ParseFloat(status["duration"], 64)
	if err != nil {
		log.Fatalln(err)
	}
	progress := int((elapsed / total) * 100)
	elapsedMin := strconv.Itoa(int(elapsed)/60) + "." + strconv.Itoa(int(elapsed)%60)
	totalMin := strconv.Itoa(int(total)/60) + "." + strconv.Itoa(int(total)%60)
	labelString := "(" + elapsedMin + "/" + totalMin + ")"
	return songstring, progress, labelString
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
	if status["volume"] == "" {
		statusString = statusString + "Volume : " + "N/A"
		return statusString
	}
	statusString = statusString + "Volume : " + status["volume"]
	return statusString
}

func (m *Music) ToggleRepeat() {
	m.Repeat = !m.Repeat
	m.Client.Repeat(m.Repeat)
}

func (m *Music) ToggleShuffle() {
	m.Shuffle = !m.Shuffle
	m.Client.Random(m.Shuffle)
}

func (m *Music) ReturnAlbums() []string {
	return []string{"Album 1", "Album 2"}
}

func (m *Music) ReturnArtists() []string {
	return []string{"Artist 1", "Artist 2"}
}

func (m *Music) LoadPlaylistIntoQueue(play string) {
	if err := m.Client.PlaylistLoad(play, -1, -1); err != nil {
		log.Fatalln(err)
	}
}

func (m *Music) GetAllSongs() []string {
	songs, err := m.Client.GetFiles()
	if err != nil {
		log.Fatalln(err)
	}
	return songs
}

func (m *Music) AddToQueue(uri string) {
	err := m.Client.Add(uri)
	if err != nil {
		log.Fatalln(err)
	}
}

func (m *Music) StopPlaying() {
	err := m.Client.Stop()
	if err != nil {
		log.Fatalln(err)
	}
}
