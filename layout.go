package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func appGrid(content tview.Primitive) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, 40, 0).
		SetRows(-1, 10, -3).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)
}

func outputPane(title string) *tview.TextView {
	tv := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	tv.SetBorder(true).SetTitle(title)
	return tv
}

func twoColumnPage(form, output tview.Primitive) *tview.Flex {
	return tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(output, 0, 1, false)
}

func styledForm() *tview.Form {
	f := tview.NewForm()
	f.SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorLightGray).
		SetButtonBackgroundColor(tcell.ColorWhite).
		SetButtonTextColor(tcell.ColorBlack)
	return f
}
