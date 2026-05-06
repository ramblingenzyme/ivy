package main

import (
	"fmt"
	"os"
	"os/exec"

	dvdrip "github.com/ramblingenzyme/ivy/dvd-rip"
	"github.com/rivo/tview"
)

func ripDVDPage(pages *tview.Pages, app *tview.Application, dvdfsMount string) tview.Primitive {
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

				session, err := dvdrip.NewSession(dvdfsMount)
				if err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}
				// defer session.Close()

				dvd, err := session.Info(device, w)
				if err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}
				fmt.Fprintf(w, "Title: %s  Main track: %s\n\n", dvd.Title, dvd.MainTrack)

				if err := session.Backup(outputDir, w); err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}

				if err := session.MergeMKV(w); err != nil {
					fmt.Fprintf(w, "Error: %v\n", err)
					return
				}

				fmt.Fprintf(w, "Done\n")
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
		// TODO: add dropdown to choose whether you're playing a real DVD, a disk ISO or a normal movie file
		AddFormItem(deviceInput).
		AddButton("Play", func() {
			device := deviceInput.GetText()
			app.Suspend(func() {
				// Have to pass dvd:// + --dvd-device, otherwise mpv doesn't skip menus/runs into other issues. This applies to DVD ISOs too.
				cmd := exec.Command(
					"cage", "-d", "--",
					"mpv", "dvd://", "--dvd-device="+device,
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
