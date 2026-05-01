package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	dvdrip "github.com/ramblingenzyme/ivy/dvd-rip"
	"github.com/rivo/tview"
)

func ripDVDPage(pages *tview.Pages, app *tview.Application) tview.Primitive {
	deviceInput := tview.NewInputField().SetLabel("Device").SetText("/dev/sr0")
	outputInput := tview.NewInputField().SetLabel("Output").SetText("/mnt/nas/movies")
	ripOutput := outputPane("rip")

	form := styledForm().
		AddFormItem(deviceInput).
		AddFormItem(outputInput).
		AddButton("Rip", func() {
			device := deviceInput.GetText()
			outputDir := outputInput.GetText()
			ripOutput.Clear()
			go func() {
				w := newOutputWriter(ripOutput, app)

				dvd, err := dvdrip.DVDInfo(device, w)
				if err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}
				fmt.Fprintf(w, "Title: %s  Main track: %s\n\n", dvd.Title, dvd.MainTrack)

				if err := dvdrip.BackupDVD(device, outputDir, w); err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}

				vocPath := filepath.Join(outputDir, dvd.Title, "VIDEO_TS", fmt.Sprintf("VTS_%s_0.VOB", dvd.MainTrack))
				mkvPath := filepath.Join(outputDir, dvd.Title+".mkv")
				if err := dvdrip.MergeMKV(vocPath, mkvPath, w); err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}

				fmt.Fprintf(w, "Done: %s\n", mkvPath)
			}()
		}).
		AddButton("Cancel", func() { pages.SwitchToPage("menu") })
	form.SetBorder(true).SetTitle("Rip DVD")

	output := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ripOutput, 0, 1, false)

	return twoColumnPage(form, output)
}

func watchDVDPage(pages *tview.Pages, app *tview.Application) tview.Primitive {
	deviceInput := tview.NewInputField().SetLabel("Device").SetText("/dev/sr0")

	form := styledForm().
		AddFormItem(deviceInput).
		AddButton("Play", func() {
			device := deviceInput.GetText()
			app.Suspend(func() {
				cmd := exec.Command("mpv", "dvd://",
					"--dvd-device="+device,
					"--vo=drm",
					"--hwdec=auto",
				)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			})
		}).
		AddButton("Cancel", func() { pages.SwitchToPage("menu") })
	form.SetBorder(true).SetTitle("Watch DVD")

	return appGrid(form)
}
