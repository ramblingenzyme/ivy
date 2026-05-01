package main

import (
	"io"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type crWriter struct {
	mu        sync.Mutex
	tv        *tview.TextView
	app       *tview.Application
	committed strings.Builder
	current   strings.Builder
}

func (w *crWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	for _, b := range p {
		switch b {
		case '\r':
			w.current.Reset()
		case '\n':
			w.committed.WriteString(w.current.String())
			w.committed.WriteByte('\n')
			w.current.Reset()
		default:
			w.current.WriteByte(b)
		}
	}
	committed := w.committed.String()
	current := w.current.String()
	w.mu.Unlock()

	w.app.QueueUpdateDraw(func() {
		w.tv.SetText(committed + current)
	})
	return len(p), nil
}

func newOutputWriter(tv *tview.TextView, app *tview.Application) io.Writer {
	return &crWriter{tv: tv, app: app}
}

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
	tv.SetChangedFunc(func() {
		tv.ScrollToEnd()
	})
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
