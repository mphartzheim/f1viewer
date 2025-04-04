package tabs

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
)

// CreateScheduleTab converts a *data.ScheduleResponse into a Fyne table widget.
// It displays columns in the order: Round, Race Name, Circuit, Location, Date.
// The "Race Name" and "Circuit" columns are rendered as hyperlinks that link to the corresponding URL.
// Modify CreateScheduleTab to accept a callback for when a flag is clicked.
func CreateScheduleTab(schedule *data.ScheduleResponse, onFlagClicked func(round string)) *widget.Table {
	races := schedule.MRData.RaceTable.Races
	rowCount := len(races) + 1 // header row + races
	colCount := 5              // Round, Race Name, Circuit, Location, Date

	// Find the index of the next upcoming race
	nextIndex := -1
	now := time.Now().UTC()
	for i, race := range races {
		if race.Date != "" && race.Time != "" {
			fullTimeStr := fmt.Sprintf("%sT%s", race.Date, race.Time)
			if t, err := time.Parse(time.RFC3339, fullTimeStr); err == nil && t.After(now) {
				nextIndex = i
				break
			}
		}
	}

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
				// Header row
				headers := []string{"Round", "Race", "Date", "Time", "Circuit", "Location"}
				updateCell(widget.NewLabel(headers[id.Col]))
				return
			}

			// Race row
			race := races[id.Row-1]
			switch id.Col {
			case 0:
				roundLabel := race.Round
				if id.Row-1 == nextIndex {
					roundLabel += " (next)"
				}
				updateCell(widget.NewLabelWithStyle(roundLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}))
			case 1:
				// Race name + finished flag logic
				u, err := url.Parse(race.URL)
				var raceHyperlink fyne.CanvasObject
				if err == nil && u != nil {
					hl := widget.NewHyperlink(race.RaceName, u)
					var extra fyne.CanvasObject
					if race.Date != "" && race.Time != "" {
						fullTimeStr := fmt.Sprintf("%sT%s", race.Date, race.Time)
						if t, err := time.Parse(time.RFC3339, fullTimeStr); err == nil && t.Before(now) {
							flagButton := widget.NewButton("🏁", func() {
								onFlagClicked(race.Round)
							})
							flagButton.Importance = widget.LowImportance
							extra = flagButton
						} else {
							extra = widget.NewLabel("\u200B")
						}
					} else {
						extra = widget.NewLabel("\u200B")
					}
					raceHyperlink = container.NewHBox(hl, extra)
				} else {
					raceHyperlink = widget.NewLabel(race.RaceName)
				}
				updateCell(raceHyperlink)
			case 2:
				// Circuit with optional map button
				u, err := url.Parse(race.Circuit.URL)
				var circuitHyperlink fyne.CanvasObject
				if err == nil && u != nil {
					circuitHyperlink = widget.NewHyperlink(race.Circuit.CircuitName, u)
				} else {
					circuitHyperlink = widget.NewLabel(race.Circuit.CircuitName)
				}

				items := []fyne.CanvasObject{circuitHyperlink}
				lat := race.Circuit.Location.Lat
				lon := race.Circuit.Location.Long
				if lat != "" && lon != "" {
					mapURLStr := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", data.MapBaseURL, lat, lon, lat, lon)
					mapURL, err := url.Parse(mapURLStr)
					if err == nil {
						mapButton := widget.NewButton("🗺️", func() {
							_ = fyne.CurrentApp().OpenURL(mapURL)
						})
						mapButton.Importance = widget.LowImportance
						items = append(items, mapButton)
					}
				}
				updateCell(container.NewHBox(items...))
			case 3:
				loc := race.Circuit.Location
				updateCell(widget.NewLabel(fmt.Sprintf("%s, %s", loc.Locality, loc.Country)))
			case 4:
				updateCell(widget.NewLabel(race.Date))
			default:
				updateCell(widget.NewLabel(""))
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 70)  // Round
	table.SetColumnWidth(1, 280) // Race
	table.SetColumnWidth(2, 280) // Circuit
	table.SetColumnWidth(3, 280) // Location
	table.SetColumnWidth(4, 120) // Date
	table.Resize(fyne.NewSize(820, float32((len(races)+1)*30)))

	return table
}
