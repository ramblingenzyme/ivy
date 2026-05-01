package main

import (
	"github.com/rivo/tview"
)

func importPhotosPage(pages *tview.Pages) tview.Primitive {
	form := styledForm().
		AddInputField("Source", "/media/sd", 0, nil, nil).
		AddInputField("Destination", "/mnt/nas/photos", 0, nil, nil).
		AddButton("Import", nil).
		AddButton("Cancel", func() { pages.SwitchToPage("menu") })
	form.SetBorder(true).SetTitle("Import Photos")

	output := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(outputPane("rsync"), 0, 1, false).
		AddItem(outputPane("exiftool"), 0, 1, false)

	return twoColumnPage(form, output)
}
