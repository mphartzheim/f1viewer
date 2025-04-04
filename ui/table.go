package ui

import (
	"fmt"

	"f1viewer/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// CreateTable creates a table widget from a slice of data.DataType.
func CreateTable(records []data.DataType) *widget.Table {
	table := widget.NewTable(
		// Returns the number of rows and columns.
		func() (int, int) {
			return len(records) + 1, 3
		},
		// Creates a new template cell.
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		// Populates each cell using the new TableCellID type.
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				// Header row.
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Name")
				case 2:
					label.SetText("Extra")
				}
			} else {
				// Data rows. Adjust index for header.
				item := records[id.Row-1]
				switch id.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", item.ID))
				case 1:
					label.SetText(item.Name)
				case 2:
					label.SetText("N/A")
				}
			}
		},
	)
	return table
}
