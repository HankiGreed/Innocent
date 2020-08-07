package tui

import (
	"time"

	"github.com/HankiGreed/Innocent/pkg/config"
	"github.com/HankiGreed/Innocent/pkg/database"
	"github.com/HankiGreed/Innocent/pkg/music"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var views = []string{" Playlists ", " Albums ", " Artists "}

type UI struct {
	Grid       *ui.Grid
	ActiveView int
	Mainview   *widgets.List
	Sideview   *widgets.List
	Songview   *widgets.Gauge
	Searchview *widgets.Paragraph
	Infoview   *widgets.Paragraph
	MPD        *music.Music
	Options    *config.Config
	Db         *database.Database
}

func (v *UI) InitializeInterface() {
	maxX, maxY := ui.TerminalDimensions()
	v.ActiveView = 0
	v.Sideview = widgets.NewList()
	v.MPD = &music.Music{}
	v.MPD.ConnectToClient()
	playlists := v.MPD.ReturnPlaylists()
	v.Sideview.Rows = playlists
	v.Sideview.Title = " Playlists "

	v.Mainview = widgets.NewList()
	v.Mainview.Title = " Songs "
	v.Mainview.Rows = v.MPD.ReturnSongsInPlaylist(playlists[0])

	v.Songview = widgets.NewGauge()
	v.Songview.Title, v.Songview.Percent, v.Songview.Label = v.MPD.GetNowPlaying()
	v.Songview.BarColor = ui.ColorBlue

	v.Searchview = widgets.NewParagraph()
	v.Searchview.Text = ""
	v.Searchview.Title = " Search "

	v.Infoview = widgets.NewParagraph()
	v.Infoview.Title = " Status "
	v.Infoview.Text = v.MPD.ReturnStatusString()

	v.Options = config.ReadConfig()

	v.Db = database.ConnectToDb(v.Options.DbConfig.Path)
	v.Grid = ui.NewGrid()
	v.Grid.SetRect(0, 0, maxX, maxY)

	v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0/5, v.Sideview), ui.NewCol(4.0/5, v.Mainview)),
		ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
}

func (v *UI) MainLoop() {

	ticker := time.NewTicker(time.Second).C
	ui.Render(v.Grid)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				v.Grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(v.Grid)
			case "<Tab>":
				v.manageSideView()
				ui.Clear()
				ui.Render(v.Grid)
			}
		case <-ticker:
			v.UpdateContent()
			ui.Clear()
			ui.Render(v.Grid)
		}
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
