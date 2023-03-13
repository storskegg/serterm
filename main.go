// Demo code for the TextView primitive.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	io2 "github.com/storskegg/serterm/io"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f:", r)
			os.Exit(100)
		}
	}()

	cmd := ""

	var sd io.ReadWriteCloser

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlX:
			app.Stop()
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		default:
			return event
		}
	})

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).
		SetTitle("[ Serial Stream ]")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		textView.SetText("Initializing serial device...")
		time.Sleep(1 * time.Second)
		textView.Clear()
		sd = io2.NewLoopBack()

		s := bufio.NewScanner(sd)
		s.Split(bufio.ScanLines)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if s.Scan() {
					line := strings.TrimSpace(s.Text())
					textView.Write([]byte(fmt.Sprintf("[orange]Response:[white] '%s'\n", line)))
				}
			}
		}
	}()

	flexCenter := tview.NewFlex().SetDirection(tview.FlexRow)
	flexCenter.AddItem(textView, 0, 1, false)

	form := tview.NewForm()

	form.AddInputField("Command", "", 20, nil, func(text string) {
		if text == "" {
			return
		}
	}).AddButton("Send", func() {
		sd.Write([]byte(strings.TrimSpace(cmd) + "\n"))
	}).AddButton("Quit", func() {
		app.Stop()
	})
	form.SetFieldBackgroundColor(tcell.ColorYellow)
	form.SetFieldTextColor(tcell.ColorBlack)
	btnStyle := tcell.Style{}.
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorWhite)
	form.SetButtonStyle(btnStyle)

	form.
		SetBorder(true).
		SetTitle("[ Commands ]")

	form.SetBackgroundColor(tcell.ColorBlack)

	flexCenter.AddItem(form, 10, 1, true)

	rootFlex := tview.NewFlex().
		AddItem(flexCenter, 0, 2, true)

	app.SetRoot(rootFlex, true)
	app.EnableMouse(false)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
