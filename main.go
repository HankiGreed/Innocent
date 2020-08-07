package main

import (
	"log"

	"github.com/HankiGreed/Innocent/pkg/tui"
	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	v := tui.UI{}
	v.InitializeInterface()
	v.MainLoop()
}
