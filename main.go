// Demo code for the TextView primitive.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
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
			debug.PrintStack()
			os.Exit(100)
		}
	}()

	cmd := ""

	sd := io2.NewLoopBack()

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

	recvWriter := io2.NewRecvPrepender(textView)
	sp := io2.NewSendPrepender(textView)
	sendWriter := io.MultiWriter(sd, sp)

	textView.SetBorder(true).
		SetTitle("[ Serial Stream ]")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		textView.SetText("Initializing serial device...")
		time.Sleep(1 * time.Second)
		textView.Clear()

		s := bufio.NewScanner(sd)
		s.Split(bufio.ScanLines)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if s.Scan() {
					recvWriter.Write([]byte(strings.TrimSpace(s.Text())))
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
		cmd = text
	}).AddButton("Send", func() {
		sendWriter.Write([]byte(strings.TrimSpace(cmd) + "\n"))
		//form.GetFormItemByLabel("Command").
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

	flexCenter.AddItem(form, 7, 1, true)

	rootFlex := tview.NewFlex().
		AddItem(flexCenter, 0, 2, true)

	app.SetRoot(rootFlex, true)
	app.EnableMouse(false)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
