package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createApp(target string, hostname string, ports []*Port) (*tview.Application, func()) {
	table := tview.NewTable()
	table.
		SetBorderPadding(0, 0, 2, 2).
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorGreen).
		SetTitle(fmt.Sprintf(" Monitoring %s (%s) ", target, hostname))

	setupHeader(table, "Port", "Alias", "RX", "TX")

	for i, port := range ports {
		row := i + 1
		table.InsertRow(row)
		table.SetCell(row, 0, tview.NewTableCell(port.Name))
		table.SetCell(row, 1, tview.NewTableCell(port.Alias))
		table.SetCell(row, 2, tview.NewTableCell(""))
		table.SetCell(row, 3, tview.NewTableCell(""))
	}

	app := tview.NewApplication().SetRoot(table, true)
	update := func() {
		for i, port := range ports {
			row := i + 1
			table.GetCell(row, 2).SetText(toReadable(port.Rx))
			table.GetCell(row, 3).SetText(toReadable(port.Tx))
		}

		app.Draw()
	}

	return app, update
}

func setupHeader(table *tview.Table, headings ...string) {
	table.InsertRow(0)
	for i := 0; i < len(headings); i++ {
		table.InsertColumn(i)
		table.SetCell(0, i, tview.NewTableCell(headings[i]).SetTextColor(tcell.ColorYellow))
	}
}

func toReadable(speed uint) string {
	const (
		K = 1000
		M = K * 1000
		G = M * 1000
	)

	if speed >= G {
		return fmt.Sprintf("%.1f GBits/s", float32(speed)/G)
	}

	if speed >= M {
		return fmt.Sprintf("%.1f Mbits/s", float32(speed)/M)
	}

	if speed >= K {
		return fmt.Sprintf("%.1f Kbits/s", float32(speed)/K)
	}

	return fmt.Sprintf("%d bits/s", speed)
}
