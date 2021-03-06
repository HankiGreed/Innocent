package tui

import (
	"strings"
	"time"

	"github.com/HankiGreed/Innocent/pkg/config"
	"github.com/HankiGreed/Innocent/pkg/music"
	ui "github.com/HankiGreed/termui/v3"
	"github.com/HankiGreed/termui/v3/widgets"
)

var views = []string{" Playlists ", " Albums ", " Artists "}

type UI struct {
	Grid                *ui.Grid
	ActiveView          int
	Queueview           *widgets.List
	AllSongsView        *widgets.Tree
	Mainview            *widgets.List
	Sideview            *widgets.List
	Songview            *widgets.Gauge
	Searchview          *widgets.Paragraph
	Infoview            *widgets.Paragraph
	MPD                 *music.Music
	Options             *config.Config
	ActivePane          string
	ActiveWindow        string
	CurrentQueueSongIDs []int
	CurrentTreePath     string
}

func (v *UI) InitializeInterface() {
	maxX, maxY := ui.TerminalDimensions()
	v.ActiveView = 0
	v.Sideview = widgets.NewList()
	v.MPD = &music.Music{}
	v.MPD.ConnectToClient()

	v.Options = config.ReadConfig()

	playlists := v.MPD.ReturnPlaylists()
	v.Sideview.Rows = playlists
	v.Sideview.Title = " Playlists "
	v.Sideview.TextStyle = ui.NewStyle(ui.ColorRed)
	v.Sideview.BorderStyle = ui.NewStyle(ui.ColorGreen)
	v.ActivePane = "side"

	v.Mainview = widgets.NewList()
	v.Mainview.Title = " Songs "
	v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(playlists[0])
	v.Mainview.TextStyle = ui.NewStyle(ui.ColorRed)

	v.Queueview = widgets.NewList()
	v.Queueview.Title = " Now Playing "
	v.Queueview.Rows, v.CurrentQueueSongIDs = v.MPD.GetCurrentQueue()
	v.Queueview.TextStyle = ui.NewStyle(ui.ColorRed)
	v.Queueview.BorderStyle = ui.NewStyle(ui.ColorGreen)

	v.AllSongsView = widgets.NewTree()
	v.AllSongsView.Title = " All Songs "
	v.AllSongsView.SetNodes(getTreeNodesFromDirectoryListing(v.Options.MusicDirectory))
	v.AllSongsView.TextStyle = ui.NewStyle(ui.ColorRed)
	v.AllSongsView.BorderStyle = ui.NewStyle(ui.ColorGreen)

	v.Songview = widgets.NewGauge()
	v.Songview.Title, v.Songview.Percent, v.Songview.Label = v.MPD.GetNowPlaying()
	v.Songview.BarColor = ui.ColorBlue

	v.Searchview = widgets.NewParagraph()
	v.Searchview.Text = ""
	v.Searchview.Title = " Search "

	v.Infoview = widgets.NewParagraph()
	v.Infoview.Title = " Status "
	v.Infoview.Text = v.MPD.ReturnStatusString()

	v.Grid = ui.NewGrid()
	v.Grid.SetRect(0, 0, maxX, maxY)
	v.ActiveWindow = "Queue"
	v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0, v.Queueview)),
		ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
}

func (v *UI) MainLoop() {

	ticker := time.NewTicker(time.Second).C
	ui.Render(v.Grid)
	uiEvents := ui.PollEvents()
	prevKey := ""
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<Escape>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				v.Grid.SetRect(0, 0, payload.Width, payload.Height)
			case "<Tab>":
				v.manageSideView()
			case "j", "<Down>":
				v.ScrollDownCurrentView()
			case "k", "<Up>":
				v.ScrollUpCurrentView()
			case "l", "<Right>":
				v.ActivePane = "main"
				v.Mainview.BorderStyle = ui.NewStyle(ui.ColorGreen)
				v.Sideview.BorderStyle = ui.NewStyle(ui.ColorClear)
			case "h", "<Left>":
				v.ActivePane = "side"
				v.Sideview.BorderStyle = ui.NewStyle(ui.ColorGreen)
				v.Mainview.BorderStyle = ui.NewStyle(ui.ColorClear)
			case "G", "End":
				v.ScrollCurrentEnd()
			case "g":
				if prevKey == "g" {
					v.ScrollCurrentStart()
				}
			case "<Home>":
				v.ScrollCurrentStart()
			case "<C-d>":
				v.ScrollCurrentHalfDown()
			case "<C-u>":
				v.ScrollCurrentHalfUp()
			case "r":
				v.MPD.ToggleRepeat()
			case "z":
				v.MPD.ToggleShuffle()
			case "<Space>":
				v.HandleSpace()
			case "a":
				v.HandleAdding()
			case "s":
				v.MPD.StopPlaying()
			case "<Enter>":
				v.HandleEnter()
			case "u":
				v.MPD.UpdateDatabase()
			}
			if prevKey == "g" {
				prevKey = ""
			} else {
				prevKey = e.ID
			}
		case <-ticker:
			v.UpdateContent()
			ui.Clear()
			ui.Render(v.Grid)
		}
		ui.Clear()
		ui.Render(v.Grid)
	}
}

func (v *UI) UpdateContent() {
	v.Infoview.Text = v.MPD.ReturnStatusString()
	v.Songview.Title, v.Songview.Percent, v.Songview.Label = v.MPD.GetNowPlaying()
	v.Songview.BarColor = ui.ColorBlue
}

func (v *UI) manageSideView() {
	v.ActiveView = (v.ActiveView + 1) % len(views)
	v.Sideview.Title = views[v.ActiveView]
	switch v.ActiveView {
	case 0:
		v.Sideview.Rows = v.MPD.ReturnPlaylists()
	case 1:
		v.Sideview.Rows = v.MPD.ReturnAlbums()
	case 2:
		v.Sideview.Rows = v.MPD.ReturnArtists()
	}
}

func (v *UI) ScrollDownCurrentView() {
	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollDown()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollDown()
		}
	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollDown()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollDown()
	}
}

func (v *UI) ScrollUpCurrentView() {

	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollUp()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollUp()
		}
	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollUp()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollUp()
	}

}

func (v *UI) ScrollCurrentEnd() {

	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollBottom()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollBottom()
		}

	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollBottom()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollBottom()
	}

}

func (v *UI) ScrollCurrentStart() {

	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollTop()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollTop()
		}

	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollTop()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollTop()
	}

}

func (v *UI) ScrollCurrentHalfDown() {

	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollHalfPageDown()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollHalfPageDown()
		}
	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollHalfPageDown()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollHalfPageDown()
	}

}

func (v *UI) ScrollCurrentHalfUp() {

	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.Sideview.ScrollHalfPageUp()
			v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(v.Sideview.Rows[v.Sideview.SelectedRow])
		} else {
			v.Mainview.ScrollHalfPageUp()
		}
	} else if v.ActiveWindow == "AllSongs" {
		v.AllSongsView.ScrollHalfPageUp()
	} else if v.ActiveWindow == "Queue" {
		v.Queueview.ScrollHalfPageUp()
	}

}

func (v *UI) HandleAdding() {
	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.MPD.LoadPlaylistIntoQueue(v.Sideview.Rows[v.Sideview.SelectedRow])
		}
	} else if v.ActiveWindow == "AllSongs" {
		v.MPD.AddToQueue(v.CurrentTreePath + "/" + v.AllSongsView.SelectedNode().Value.String())
		v.AllSongsView.ScrollDown()
	}
}

func (v *UI) HandleSpace() {
	maxX, maxY := ui.TerminalDimensions()
	if v.ActiveWindow == "AllSongs" {
		v.ActiveWindow = "Home"
		v.Grid = ui.NewGrid()
		v.Grid.SetRect(0, 0, maxX, maxY)
		v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0/5, v.Sideview), ui.NewCol(4.0/5, v.Mainview)),
			ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
	} else if v.ActiveWindow == "Home" {
		v.ActiveWindow = "Queue"
		v.Grid = ui.NewGrid()
		v.Queueview.Rows, v.CurrentQueueSongIDs = v.MPD.GetCurrentQueue()
		v.Grid.SetRect(0, 0, maxX, maxY)
		v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0, v.Queueview)),
			ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
	} else if v.ActiveWindow == "Queue" {
		v.ActiveWindow = "AllSongs"
		v.Grid = ui.NewGrid()
		v.Grid.SetRect(0, 0, maxX, maxY)
		v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0, v.AllSongsView)),
			ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
	}
}

func (v *UI) HandleEnter() {
	if v.ActiveWindow == "Home" {
		if v.ActivePane == "side" {
			v.MPD.Client.Clear()
			v.MPD.LoadPlaylistIntoQueue(v.Sideview.Rows[v.Sideview.SelectedRow])
			v.MPD.Client.Play(-1)
		}
	} else if v.ActiveWindow == "AllSongs" {

		if !v.AllSongsView.SelectedNode().Expanded {
			if v.AllSongsView.SelectedNode().Level == 0 {
				v.CurrentTreePath = ""
			}
			if v.CurrentTreePath == "" {
				v.CurrentTreePath += v.AllSongsView.SelectedNode().Value.String()
			} else {
				v.CurrentTreePath += ("/" + v.AllSongsView.SelectedNode().Value.String())
			}
			v.AllSongsView.SelectedNode().Nodes = getTreeNodesFromDirectoryListing(v.Options.MusicDirectory + "/" + v.CurrentTreePath)
			v.AllSongsView.ToggleExpand()
		} else {
			if v.CurrentTreePath != "" {
				if ind := strings.LastIndex(v.CurrentTreePath, "/"); ind != -1 {
					if ind == len(v.CurrentTreePath)-1 {
						v.CurrentTreePath = ""
					}
					v.CurrentTreePath = v.CurrentTreePath[:ind]
				} else {
					v.CurrentTreePath = ""
				}
			}
			v.AllSongsView.ToggleExpand()
		}
	} else if v.ActiveWindow == "Queue" {
		v.MPD.PlayID(v.CurrentQueueSongIDs[v.Queueview.SelectedRow])
	}
}
