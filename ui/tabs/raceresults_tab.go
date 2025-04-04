package tabs

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
)

// createRaceResultsTab converts a *data.RaceResultsResponse into a Fyne table widget.
// It flattens the results so that each row corresponds to one driver's result.
// Columns: Race, Round, Circuit, Date, Driver, Finish Pos, Grid, Laps, Points, Status.
func CreateRaceResultsTab(results *data.RaceResultsResponse) *widget.Table {
	// Define a helper struct to hold the flattened row data.
	type resultRow struct {
		Position       string // finishing position
		Driver         string // driver's full name
		DriverURL      string
		Points         string
		Constructor    string // constructor name
		ConstructorURL string
		Grid           string
		Laps           string
		TimeOrStatus   string
	}
	var rows []resultRow

	// Flatten the race results: one row per driver result.
	for _, race := range results.MRData.RaceTable.Races {
		for _, res := range race.Results {
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
				Points:         res.Points,
				Constructor:    res.Constructor.Name,
				ConstructorURL: res.Constructor.URL,
				Grid:           res.Grid,
				Laps:           res.Laps,
				TimeOrStatus:   timeOrStatus,
			}
			rows = append(rows, row)
		}
	}

	rowCount := len(rows) + 1 // add one for header row
	colCount := 7             // Finish Pos, Driver, Points, Constructor, Grid, Laps, Time/Status

	table := widget.NewTable(
		func() (int, int) { return rowCount, colCount },
		// Template cell: use a Stack container so we can swap in a container when needed.
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// In this implementation, we'll check if our cell is a container.
			cont, ok := cell.(*fyne.Container)
			if !ok {
				return
			}

			// Helper: update the container with a new CanvasObject.
			updateCell := func(obj fyne.CanvasObject) {
				cont.Objects = []fyne.CanvasObject{obj}
				cont.Refresh()
			}

			if id.Row == 0 {
				// Header row.
				var header string
				switch id.Col {
				case 0:
					header = "Pos"
				case 1:
					header = "Driver"
				case 2:
					header = "Points"
				case 3:
					header = "Constructor"
				case 4:
					header = "Grid"
				case 5:
					header = "Laps"
				case 6:
					header = "Time/Status"
				default:
					header = ""
				}
				updateCell(widget.NewLabel(header))
				return
			}

			// Data rows: get the corresponding row (adjust for header).
			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				updateCell(widget.NewLabel(row.Position))
			case 1:
				{
					// Driver column: create a hyperlink using the driver's URL.
					u, err := url.Parse(row.DriverURL)
					var driverObj fyne.CanvasObject
					if err == nil && u != nil {
						hl := widget.NewHyperlink(row.Driver, u)
						hl.Resize(hl.MinSize())
						// Append an invisible label.
						inv := widget.NewLabel("\u200B")
						inv.Resize(inv.MinSize())
						driverObj = container.NewHBox(hl, inv)
					} else {
						driverObj = widget.NewLabel(row.Driver)
					}
					updateCell(driverObj)
				}
			case 2:
				updateCell(widget.NewLabel(row.Points))
			case 3:
				{
					// Constructor column: create a hyperlink using the constructor's URL.
					u, err := url.Parse(row.ConstructorURL)
					var consObj fyne.CanvasObject
					if err == nil && u != nil {
						hl := widget.NewHyperlink(row.Constructor, u)
						hl.Resize(hl.MinSize())
						// Append an invisible label.
						inv := widget.NewLabel("\u200B")
						inv.Resize(inv.MinSize())
						consObj = container.NewHBox(hl, inv)
					} else {
						consObj = widget.NewLabel(row.Constructor)
					}
					updateCell(consObj)
				}
			case 4:
				updateCell(widget.NewLabel(row.Grid))
			case 5:
				updateCell(widget.NewLabel(row.Laps))
			case 6:
				updateCell(widget.NewLabel(row.TimeOrStatus))
			default:
				updateCell(widget.NewLabel(""))
			}
		},
	)

	// Set desired column widths.
	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 200) // Driver
	table.SetColumnWidth(2, 70)  // Points
	table.SetColumnWidth(3, 200) // Constructor
	table.SetColumnWidth(4, 70)  // Grid
	table.SetColumnWidth(5, 70)  // Laps
	table.SetColumnWidth(6, 120) // Time/Status

	// Resize the table based on total width and a row height of 30 pixels.
	table.Resize(fyne.NewSize(950, float32(rowCount*30)))
	return table
}
