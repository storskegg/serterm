// Demo code for the TextView primitive.
package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const corporate = `Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.

Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.

Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.

[yellow]Press Enter, then Tab/Backtab for word selections`

func streamer(w io.Writer, numSelections *int) {
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

	numSelections := 0
	go streamer(textView, &numSelections)

	textView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := textView.GetHighlights()

		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				textView.Highlight()
			} else {
				textView.Highlight("0").ScrollToHighlight()
			}
		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab {
				index = (index + 1) % numSelections
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections
			} else {
				return
			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})

	textView.SetBorder(true).
		SetTitle("Middle (3 x height of Top)")

	flexCenter := tview.NewFlex().SetDirection(tview.FlexRow)
	flexCenter.AddItem(textView, 0, 1, true)

	boxCenterBottom := tview.NewBox().
		SetBorder(true).
		SetTitle("Bottom (5 rows)").
		SetBackgroundColor(tcell.ColorBlack)
	flexCenter.AddItem(boxCenterBottom, 5, 1, false)

	rootFlex := tview.NewFlex().
		AddItem(flexCenter, 0, 2, true)
	
	if err := app.SetRoot(rootFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
