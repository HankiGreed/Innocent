package tui

import (
	"github.com/HankiGreed/Innocent/pkg/music"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type UI struct {
	Grid       *ui.Grid
	ActiveView int
	Mainview   *widgets.List
	Sideview   *widgets.List
	Songview   *widgets.Gauge
	Searchview *widgets.Paragraph
	Infoview   *widgets.Paragraph
}

func (v *UI) InitializeGrid() {
	maxX, maxY := ui.TerminalDimensions()

	v.Sideview = widgets.NewList()
	mpd := music.Music{}
	mpd.ConnectToClient()
	playlists := mpd.ReturnPlaylists()
	v.Sideview.Rows = playlists
	v.Sideview.Title = " Playlists "

	v.Mainview = widgets.NewList()
	v.Mainview.Title = " Songs "
	v.Mainview.Rows = mpd.ReturnSongsInPlaylist(playlists[0])

	v.Songview = widgets.NewGauge()
	v.Songview.Title = mpd.GetNowPlaying()
	v.Songview.Percent = 50
	v.Songview.BarColor = ui.ColorBlue

	v.Searchview = widgets.NewParagraph()
	v.Searchview.Text = ""
	v.Searchview.Title = " Search "

	v.Infoview = widgets.NewParagraph()
	v.Infoview.Title = " Status "
	v.Infoview.Text = mpd.ReturnStatusString()

	v.Grid = ui.NewGrid()
	v.Grid.SetRect(0, 0, maxX, maxY)

	v.Grid.Set(ui.NewRow(0.6/6, ui.NewCol(2.0/3, v.Searchview), ui.NewCol(1.0/3, v.Infoview)), ui.NewRow(4.8/6, ui.NewCol(1.0/5, v.Sideview), ui.NewCol(4.0/5, v.Mainview)),
		ui.NewRow(0.6/6, ui.NewCol(1, v.Songview)))
}

func (v *UI) MainLoop() {

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
			}
		}
	}
}
