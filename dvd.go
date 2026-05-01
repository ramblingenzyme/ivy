package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ripDVDPage(pages *tview.Pages) tview.Primitive {
	format := tview.NewDropDown().
		SetLabel("Format ").
		SetOptions([]string{"MKV", "MP4", "ISO"}, nil).
		SetCurrentOption(0)
	format.SetListStyles(
		tcell.StyleDefault.Background(tcell.ColorDarkSlateGray).Foreground(tcell.ColorWhite),
		tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack),
	)

	form := styledForm().
		AddInputField("Device", "/dev/sr0", 0, nil, nil).
		AddInputField("Output", "/mnt/nas/movies", 0, nil, nil).
		AddFormItem(format).
		AddButton("Rip", nil).
		AddButton("Cancel", func() { pages.SwitchToPage("menu") })
	form.SetBorder(true).SetTitle("Rip DVD")

	output := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(outputPane("HandBrakeCLI"), 0, 1, false)

	return twoColumnPage(form, output)
}
