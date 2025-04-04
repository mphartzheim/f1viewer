package tabs

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
)

func CreateConstructorStandingsTab(results *data.ConstructorStandingsResponse) *widget.Table {
	type rowData struct {
		Position    string
		Name        string
		URL         string
		Nationality string
		Points      string
		Wins        string
	}

	var rows []rowData

	for _, standingsList := range results.MRData.StandingsTable.StandingsLists {
		for _, standing := range standingsList.ConstructorStandings {
			c := standing.Constructor
			rows = append(rows, rowData{
				Position:    standing.Position,
				Name:        c.Name,
				URL:         c.URL,
				Nationality: c.Nationality,
				Points:      standing.Points,
				Wins:        standing.Wins,
			})
		}
	}

	rowCount := len(rows) + 1
	colCount := 5 // Pos, Name, Nationality, Points, Wins

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
				headers := []string{"Pos", "Constructor", "Nationality", "Points", "Wins"}
				update(widget.NewLabel(headers[id.Col]))
				return
			}

			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				update(widget.NewLabel(row.Position))
			case 1:
				u, err := url.Parse(row.URL)
				if err == nil && u != nil {
					update(container.NewHBox(widget.NewHyperlink(row.Name, u), widget.NewLabel("\u200B")))
				} else {
					update(widget.NewLabel(row.Name))
				}
			case 2:
				update(widget.NewLabel(row.Nationality))
			case 3:
				update(widget.NewLabel(row.Points))
			case 4:
				update(widget.NewLabel(row.Wins))
			default:
				update(widget.NewLabel(""))
			}
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 200) // Constructor
	table.SetColumnWidth(2, 120) // Nationality
	table.SetColumnWidth(3, 80)  // Points
	table.SetColumnWidth(4, 60)  // Wins

	table.Resize(fyne.NewSize(750, float32(rowCount*30)))
	return table
}
