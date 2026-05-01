package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	list := tview.NewList().
		AddItem("Import photos", "Copy photos from SD card to NAS", '1', func() {
			pages.SwitchToPage("import")
		}).
		AddItem("Rip DVD", "", '2', func() {
			pages.SwitchToPage("rip")
		}).
		AddItem("Watch DVD", "", '3', func() {
			pages.SwitchToPage("watch")
		}).
		AddItem("Manage ebooks", "", '4', func() {
			pages.SwitchToPage("ebooks")
		})

	menuFrame := tview.NewFrame(list).
		AddText("OPTIONS", true, tview.AlignLeft, tcell.ColorWhite)

	pages.
		AddPage("menu", appGrid(menuFrame), true, true).
		AddPage("import", importPhotosPage(pages), true, false).
		AddPage("rip", ripDVDPage(pages, app), true, false).
		AddPage("watch", watchDVDPage(pages, app), true, false).
		AddPage("ebooks", tview.NewTextView().SetText("Manage ebooks — coming soon"), true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}

		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'q':
			name, _ := pages.GetFrontPage()
			if name != "menu" {
				pages.SwitchToPage("menu")
				return nil
			}
		}

		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
