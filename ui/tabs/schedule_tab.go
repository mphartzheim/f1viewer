package tabs

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
	"github.com/mphartzheim/f1viewer/util"
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
				headers := []string{"Round", "Race", "Circuit", "Location", "Date"}
				updateCell(widget.NewLabel(headers[id.Col]))
				return
			}

			// Race row
			race := races[id.Row-1]
			isNext := id.Row-1 == nextIndex
			primaryColor := theme.PrimaryColor()

			switch id.Col {
			case 0:
				roundLabel := race.Round
				if isNext {
					roundLabel = "Next"
				}
				// Use default foreground unless next race.
				col := theme.ForegroundColor()
				if isNext {
					col = primaryColor
				}
				ct := util.NewColoredText(roundLabel, col)
				ct.Text.TextStyle = fyne.TextStyle{}
				ct.Text.Alignment = fyne.TextAlignLeading
				updateCell(ct)

			case 1:
				u, err := url.Parse(race.URL)
				var raceCell fyne.CanvasObject
				if err == nil && u != nil {
					hl := widget.NewHyperlink(race.RaceName, u) // leave hyperlink styling intact
					var extra fyne.CanvasObject
					if isNext {
						extra = widget.NewButton("Spoilers", func() {
							onFlagClicked(race.Round)
						})
						extra.(*widget.Button).Importance = widget.LowImportance
					} else if race.Date != "" && race.Time != "" {
						fullTimeStr := fmt.Sprintf("%sT%s", race.Date, race.Time)
						if t, err := time.Parse(time.RFC3339, fullTimeStr); err == nil && t.Before(now) {
							extra = widget.NewButton("🏁", func() {
								onFlagClicked(race.Round)
							})
							extra.(*widget.Button).Importance = widget.LowImportance
						} else {
							extra = widget.NewLabel("\u200B")
						}
					} else {
						extra = widget.NewLabel("\u200B")
					}
					raceCell = container.NewHBox(hl, extra)
				} else {
					// Fallback: plain text with optional Spoilers button for next race.
					var txt fyne.CanvasObject
					if isNext {
						txt = util.NewColoredText(race.RaceName, primaryColor)
					} else {
						txt = util.NewColoredText(race.RaceName, theme.ForegroundColor())
					}
					txt.(*util.ColoredText).Text.Alignment = fyne.TextAlignLeading

					if isNext {
						spoilersButton := widget.NewButton("Spoilers", func() {
							onFlagClicked(race.Round)
						})
						spoilersButton.Importance = widget.LowImportance
						raceCell = container.NewHBox(txt, spoilersButton)
					} else {
						raceCell = txt
					}
				}
				updateCell(raceCell)

			case 2:
				u, err := url.Parse(race.Circuit.URL)
				var circuitCell fyne.CanvasObject
				if err == nil && u != nil {
					circuitCell = widget.NewHyperlink(race.Circuit.CircuitName, u) // leave hyperlink styling intact
				} else {
					// Use coloredText for plain text
					col := theme.ForegroundColor()
					if isNext {
						col = primaryColor
					}
					ct := util.NewColoredText(race.Circuit.CircuitName, col)
					ct.Text.Alignment = fyne.TextAlignLeading
					circuitCell = ct
				}
				items := []fyne.CanvasObject{circuitCell}
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
				label := fmt.Sprintf("%s, %s", loc.Locality, loc.Country)
				col := theme.ForegroundColor()
				if isNext {
					col = primaryColor
				}
				ct := util.NewColoredText(label, col)
				ct.Text.Alignment = fyne.TextAlignLeading
				updateCell(ct)

			case 4:
				col := theme.ForegroundColor()
				if isNext {
					col = primaryColor
				}
				ct := util.NewColoredText(race.Date, col)
				ct.Text.Alignment = fyne.TextAlignLeading
				updateCell(ct)

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
