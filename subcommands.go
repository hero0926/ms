package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func status() {
	var q = NewQueue()
	err := q.Load()
	if err != nil {
		if q.Save() != nil {
			fmt.Println(err)
		}
		fmt.Println(err)
		os.Exit(1)
	}

	q.ShowCurrent()
}

func newblock() {
	var (
		scanner = bufio.NewScanner(os.Stdin)
		q       = NewQueue()
	)

	err := q.Load()
	if err != nil {
		fmt.Printf("Error reading queue: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Name for new block: ")
	scanner.Scan()
	name := scanner.Text()

records:
	fmt.Printf("Amount of records for this week: ")
	scanner.Scan()
	num := scanner.Text()

	n, err := strconv.Atoi(num)
	if n <= 0 {
		err = fmt.Errorf("amount of records must be an integer greater than 0")
	}
	if err != nil {
		fmt.Printf("Error with amount of record: %s\n", err)
		fmt.Println("Pleae try again.")
		goto records
	}

	albums := make([]Album, n)

	for i := 0; i < n; i++ {
		fmt.Printf("Album %d: ", i+1)
		scanner.Scan()
		al := scanner.Text()
		if al == "" {
			break
		}

		albums[i].Name = al
	}

	b := NewBlock(name, albums...)
	q.Add(b)

	err = q.Save()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmdListen() {
	var q = NewQueue()

	if err := q.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TUI
	var (
		app  = tview.NewApplication()
		flex = tview.NewFlex()
	)

	flex.SetDirection(tview.FlexRow)

	// Borders
	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft

	// Queue
	for i, block := range q {
		// Skip resting week
		if i == 2 {
			continue
		}

		// Each block
		blocktable := tview.NewTable()
		blocktable.
			SetSelectable(true, true).
			SetFixed(1, 3).
			SetTitle("[::b] " + block.Name + " ").
			SetTitleAlign(tview.AlignLeft).
			SetTitleColor(tcell.ColorGreen).
			SetBorder(true)

		// Headers
		albumcell := tview.NewTableCell("Albums")
		albumcell.
			SetSelectable(false).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold)

		listenedcell := tview.NewTableCell("Listened")
		listenedcell.
			SetSelectable(false).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold)

		ratedcell := tview.NewTableCell("Rated")
		ratedcell.
			SetSelectable(false).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold)

		// Add headers
		blocktable.SetCell(0, 0, albumcell)
		blocktable.SetCell(0, 1, listenedcell)
		blocktable.SetCell(0, 2, ratedcell)

		// Albums
		for j, album := range block.Albums {
			var listened, rated = "❌", ""

			cell := tview.NewTableCell(album.Name)
			cell.
				SetExpansion(1).
				SetMaxWidth(len(album.Name))

			blocktable.SetCell(j+1, 0, cell)

			// Check per block
			if i == 0 && album.FirstListen {
				listened = "✓"
			}
			if i == 1 && album.SecondListen {
				listened = "✓"
			}
			if i == 3 && album.ThirdListen {
				listened = "✓"
			}
			if album.Rated {
				rated = "✓"
			}

			blocktable.SetCell(j+1, 1, tview.NewTableCell(listened).
				SetAlign(tview.AlignCenter))
			blocktable.SetCell(j+1, 2, tview.NewTableCell(rated).
				SetAlign(tview.AlignCenter))
		}

		flex.AddItem(blocktable, 0, 1, true)
	}

	// Run the App
	app.SetRoot(flex, true).SetFocus(flex)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
