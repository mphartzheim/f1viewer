package tabs

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
)

// CreateQualifyingResultsTab builds a table from QualifyingResultsResponse
func CreateQualifyingResultsTab(results *data.QualifyingResponse) *widget.Table {
	type resultRow struct {
		Position       string
		Driver         string
		DriverURL      string
		Constructor    string
		ConstructorURL string
		Q1             string
		Q2             string
		Q3             string
	}
	var rows []resultRow

	for _, race := range results.MRData.RaceTable.Races {
		for _, res := range race.QualifyingResults {
			row := resultRow{
				Position:       res.Position,
				Driver:         fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
				DriverURL:      res.Driver.URL,
				Constructor:    res.Constructor.Name,
				ConstructorURL: res.Constructor.URL,
				Q1:             res.Q1,
				Q2:             res.Q2,
				Q3:             res.Q3,
			}
			rows = append(rows, row)
		}
	}

	rowCount := len(rows) + 1
	colCount := 6 // Pos, Driver, Constructor, Q1, Q2, Q3

	table := widget.NewTable(
		func() (int, int) { return rowCount, colCount },
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			cont, ok := cell.(*fyne.Container)
			if !ok {
				return
			}
			updateCell := func(obj fyne.CanvasObject) {
				cont.Objects = []fyne.CanvasObject{obj}
				cont.Refresh()
			}
			if id.Row == 0 {
				var header string
				switch id.Col {
				case 0:
					header = "Pos"
				case 1:
					header = "Driver"
				case 2:
					header = "Constructor"
				case 3:
					header = "Q1"
				case 4:
					header = "Q2"
				case 5:
					header = "Q3"
				default:
					header = ""
				}
				updateCell(widget.NewLabel(header))
				return
			}

			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				updateCell(widget.NewLabel(row.Position))
			case 1:
				u, err := url.Parse(row.DriverURL)
				var obj fyne.CanvasObject
				if err == nil && u != nil {
					obj = container.NewHBox(widget.NewHyperlink(row.Driver, u), widget.NewLabel("\u200B"))
				} else {
					obj = widget.NewLabel(row.Driver)
				}
				updateCell(obj)
			case 2:
				u, err := url.Parse(row.ConstructorURL)
				var obj fyne.CanvasObject
				if err == nil && u != nil {
					obj = container.NewHBox(widget.NewHyperlink(row.Constructor, u), widget.NewLabel("\u200B"))
				} else {
					obj = widget.NewLabel(row.Constructor)
				}
				updateCell(obj)
			case 3:
				updateCell(widget.NewLabel(row.Q1))
			case 4:
				updateCell(widget.NewLabel(row.Q2))
			case 5:
				updateCell(widget.NewLabel(row.Q3))
			default:
				updateCell(widget.NewLabel(""))
			}
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 200) // Driver
	table.SetColumnWidth(2, 200) // Constructor
	table.SetColumnWidth(3, 90)  // Q1
	table.SetColumnWidth(4, 90)  // Q2
	table.SetColumnWidth(5, 90)  // Q3

	table.Resize(fyne.NewSize(830, float32(rowCount*30)))
	return table
}
