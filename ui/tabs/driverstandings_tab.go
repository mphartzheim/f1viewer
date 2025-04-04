package tabs

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
)

// CreateDriverStandingsTab builds a table from DriverStandingsResponse
func CreateDriverStandingsTab(results *data.DriverStandingsResponse) *widget.Table {
	type rowData struct {
		Position       string
		Driver         string
		DriverURL      string
		Nationality    string
		Constructor    string
		ConstructorURL string
		Points         string
		Wins           string
	}
	var rows []rowData

	for _, standingsList := range results.MRData.StandingsTable.StandingsLists {
		for _, standing := range standingsList.DriverStandings {
			driver := standing.Driver
			constructorName := "—"
			if len(standing.Constructors) > 0 {
				constructorName = standing.Constructors[0].Name
			}
			rows = append(rows, rowData{
				Position:       standing.Position,
				Driver:         fmt.Sprintf("%s %s", driver.GivenName, driver.FamilyName),
				DriverURL:      driver.URL,
				Nationality:    driver.Nationality,
				Constructor:    constructorName,
				ConstructorURL: standing.Constructors[0].URL,
				Points:         standing.Points,
				Wins:           standing.Wins,
			})
		}
	}

	rowCount := len(rows) + 1 // +1 for header
	colCount := 6             // Pos, Driver, Nationality, Constructor, Points, Wins

	table := widget.NewTable(
		func() (int, int) { return rowCount, colCount },
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			cont := cell.(*fyne.Container)
			update := func(obj fyne.CanvasObject) {
				cont.Objects = []fyne.CanvasObject{obj}
				cont.Refresh()
			}

			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Nationality", "Constructor", "Points", "Wins"}
				update(widget.NewLabel(headers[id.Col]))
				return
			}

			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				update(widget.NewLabel(row.Position))
			case 1:
				u, err := url.Parse(row.DriverURL)
				if err == nil && u != nil {
					update(container.NewHBox(widget.NewHyperlink(row.Driver, u), widget.NewLabel("\u200B")))
				} else {
					update(widget.NewLabel(row.Driver))
				}
			case 2:
				update(widget.NewLabel(row.Nationality))
			case 3:
				u, err := url.Parse(row.ConstructorURL)
				if err == nil && u != nil {
					update(container.NewHBox(widget.NewHyperlink(row.Constructor, u), widget.NewLabel("\u200B")))
				} else {
					update(widget.NewLabel(row.Constructor))
				}

			case 4:
				update(widget.NewLabel(row.Points))
			case 5:
				update(widget.NewLabel(row.Wins))
			default:
				update(widget.NewLabel(""))
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 180) // Driver
	table.SetColumnWidth(2, 120) // Nationality
	table.SetColumnWidth(3, 180) // Constructor
	table.SetColumnWidth(4, 80)  // Points
	table.SetColumnWidth(5, 60)  // Wins

	table.Resize(fyne.NewSize(800, float32(rowCount*30)))
	return table
}
