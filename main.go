// Demo code for the TextView primitive.
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const corporate = `Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.

Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.

Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.

[yellow]Press Enter, then Tab/Backtab for word selections[white]`

func streamer(w *tview.TextView, numSelections *int) {
	fmt.Fprintln(w, "Opening serial device...")

	time.Sleep(8 * time.Second)

	w.Clear()

	for _, word := range strings.Split(corporate, " ") {
		if word == "the" {
			word = "[#ff0000]the[white]"
		}
		if word == "to" {
			word = fmt.Sprintf(`["%d"]to[""]`, *numSelections)
			*numSelections++
		}

		fmt.Fprintf(w, "%s ", word)
		time.Sleep(25 * time.Millisecond)
	}

}

var port = 0

func main() {
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

	numSelections := 0
	go streamer(textView, &numSelections)

	flexCenter := tview.NewFlex().SetDirection(tview.FlexRow)
	flexCenter.AddItem(textView, 0, 1, false)

	form := tview.NewForm()

	form.AddInputField("Port", "0", 5, func(text string, lastChar rune) bool {
		i, err := strconv.Atoi(text)
		if err != nil {
			return false
		}
		return i >= 0 && i < 24
	}, func(text string) {
		if text == "" {
			return
		}

		p, err := strconv.Atoi(text)
		if err == nil {
			port = p
		}
	}).AddButton("Send", func() {
		fmt.Fprintf(textView, "changed: %d\n", port)
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
