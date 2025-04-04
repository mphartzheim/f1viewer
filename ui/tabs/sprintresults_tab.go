package tabs

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
)

// CreateSprintResultsTab builds a table from SprintResultsResponse
func CreateSprintResultsTab(results *data.SprintResultsResponse) *widget.Table {
	type resultRow struct {
		Position       string
		Driver         string
		DriverURL      string
		Constructor    string
		ConstructorURL string
		Points         string
		TimeOrStatus   string
	}
	var rows []resultRow

	for _, race := range results.MRData.RaceTable.Races {
		for _, res := range race.SprintResults {
			var timeOrStatus string
			if res.Time != nil && res.Time.Time != "" {
				timeOrStatus = res.Time.Time
			} else {
				timeOrStatus = res.Status
			}
			row := resultRow{
				Position:       res.Position,
				Driver:         fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
				DriverURL:      res.Driver.URL,
				Constructor:    res.Constructor.Name,
				ConstructorURL: res.Constructor.URL,
				Points:         res.Points,
				TimeOrStatus:   timeOrStatus,
			}
			rows = append(rows, row)
		}
	}

	rowCount := len(rows) + 1
	colCount := 6 // Pos, Driver, Constructor, Points, Time/Status

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
				headers := []string{"Pos", "Driver", "Constructor", "Points", "Time/Status", ""}
				if id.Col < len(headers) {
					updateCell(widget.NewLabel(headers[id.Col]))
				} else {
					updateCell(widget.NewLabel(""))
				}
				return
			}

			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				updateCell(widget.NewLabel(row.Position))
			case 1:
				u, err := url.Parse(row.DriverURL)
				if err == nil && u != nil {
					updateCell(container.NewHBox(widget.NewHyperlink(row.Driver, u), widget.NewLabel("\u200B")))
				} else {
					updateCell(widget.NewLabel(row.Driver))
				}
			case 2:
				u, err := url.Parse(row.ConstructorURL)
				if err == nil && u != nil {
					updateCell(container.NewHBox(widget.NewHyperlink(row.Constructor, u), widget.NewLabel("\u200B")))
				} else {
					updateCell(widget.NewLabel(row.Constructor))
				}
			case 3:
				updateCell(widget.NewLabel(row.Points))
			case 4:
				updateCell(widget.NewLabel(row.TimeOrStatus))
			default:
				updateCell(widget.NewLabel(""))
			}
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 200) // Driver
	table.SetColumnWidth(2, 200) // Constructor
	table.SetColumnWidth(3, 70)  // Points
	table.SetColumnWidth(4, 120) // Time/Status

	table.Resize(fyne.NewSize(800, float32(rowCount*30)))
	return table
}
