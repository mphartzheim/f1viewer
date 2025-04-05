package tabs

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
	"github.com/mphartzheim/f1viewer/userprefs"
	"github.com/mphartzheim/f1viewer/util"
)

// CreateUpcomingTab converts an UpcomingResponse into a Fyne table widget.
// It displays a header row with "Session", "Date", and "Time" followed by rows for each session.
// The "Sprint" and "Sprint Qualifying" rows are only added if data is available.
// Debugging output prints the contents of each row.
func CreateUpcomingTab(upcoming *data.UpcomingResponse) *widget.Table {
	if len(upcoming.MRData.RaceTable.Races) == 0 {
		table := widget.NewTable(
			func() (int, int) { return 1, 1 },
			func() fyne.CanvasObject { return container.NewStack(widget.NewLabel("")) },
			func(id widget.TableCellID, cell fyne.CanvasObject) {
				if cont, ok := cell.(*fyne.Container); ok {
					cont.Objects = []fyne.CanvasObject{widget.NewLabel("No upcoming race data available")}
					cont.Refresh()
				}
			},
		)
		table.SetColumnWidth(0, 360)
		table.Resize(fyne.NewSize(360, 30))
		return table
	}

	race := upcoming.MRData.RaceTable.Races[0]

	type sessionRow struct {
		Session string
		Date    string
		Time    string
	}

	rows := []sessionRow{
		{"Practice 1", race.Practice1.Date, race.Practice1.Time},
		{"Practice 2", race.Practice2.Date, race.Practice2.Time},
		{"Practice 3", race.Practice3.Date, race.Practice3.Time},
		{"Qualifying", race.Qualifying.Date, race.Qualifying.Time},
	}

	if race.Sprint.Date != "" {
		rows = append(rows, sessionRow{"Sprint", race.Sprint.Date, race.Sprint.Time})
	}

	rows = append(rows, sessionRow{"Race", race.Date, race.Time})

	// Get track lat/lon
	latStr := race.Circuit.Location.Lat
	lonStr := race.Circuit.Location.Long

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Println("Invalid latitude:", latStr)
		lat = 0
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		fmt.Println("Invalid longitude:", lonStr)
		lon = 0
	}

	// Fetch forecast data
	var forecasts []string
	for _, row := range rows {
		fullTimeStr := fmt.Sprintf("%sT%s", row.Date, row.Time)
		t, err := time.Parse(time.RFC3339, fullTimeStr)
		if err != nil {
			forecasts = append(forecasts, "Invalid time")
			continue
		}

		forecast, err := data.GetForecastForTime(lat, lon, t)
		if err != nil {
			forecasts = append(forecasts, "No data")
		} else {
			forecasts = append(forecasts, forecast)
		}
	}
	zoom := 15 // example zoom level (range is roughly 0–20; higher is more zoomed in)
	forecastURL, _ := url.Parse(fmt.Sprintf(data.WindyBaseURL, lat, lon, zoom))

	rowCount := len(rows) + 1
	colCount := 4

	var showLiveBtn binding.Bool = binding.NewBool()
	showLiveBtn.Set(false)

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

			// Header row
			if id.Row == 0 {
				headers := []string{"Session", "Forecast", "Date", "Time"}
				if id.Col < len(headers) {
					cont.Objects = []fyne.CanvasObject{widget.NewLabel(headers[id.Col])}
					cont.Refresh()
				}
				return
			}

			// Row safety check
			if id.Row-1 >= len(rows) {
				cont.Objects = []fyne.CanvasObject{widget.NewLabel("N/A")}
				cont.Refresh()
				return
			}

			row := rows[id.Row-1]
			var forecast string
			if id.Row-1 < len(forecasts) {
				forecast = forecasts[id.Row-1]
			} else {
				forecast = "N/A"
			}

			fullTimeStr := fmt.Sprintf("%sT%s", row.Date, row.Time)
			t, err := time.Parse(time.RFC3339, fullTimeStr)
			useLocal, _ := userprefs.Get().UseLocalTime.Get()
			localTime := t.Local()
			localDateStr := row.Date
			localTimeStr := row.Time
			if err == nil {
				if useLocal {
					localDateStr = localTime.Format("2006-01-02")
					localTimeStr = util.FormatTime(localTime)
				} else {
					localDateStr = t.Format("2006-01-02")
					localTimeStr = util.FormatTime(t)
				}
			}

			var display fyne.CanvasObject

			switch id.Col {
			case 0:
				display = widget.NewLabel(row.Session)
			case 1:
				forecastBox := container.NewHBox()
				btn := widget.NewButton(forecast, func() {
					_ = fyne.CurrentApp().OpenURL(forecastURL)
				})
				btn.Importance = widget.HighImportance
				forecastBox.Add(btn)
				display = forecastBox
			case 2:
				display = widget.NewLabel(localDateStr)
			case 3:
				timeLabel := widget.NewLabel(localTimeStr)
				liveButtonContainer := container.NewHBox()

				// Only one listener per row (no memory leak risk)
				go func() {
					showLiveBtn.AddListener(binding.NewDataListener(func() {
						liveButtonContainer.Objects = nil
						if v, err := showLiveBtn.Get(); err == nil && v {
							if isSessionActive(row.Date, row.Time, row.Session, time.Now()) {
								if u, err := url.Parse(data.F1tvURL); err == nil {
									btn := widget.NewButton("🔴 Live", func() {
										_ = fyne.CurrentApp().OpenURL(u)
									})
									btn.Importance = widget.HighImportance
									liveButtonContainer.Objects = append(liveButtonContainer.Objects, btn)
								}
							}
						}
						liveButtonContainer.Refresh()
					}))
				}()

				display = container.NewHBox(timeLabel, layout.NewSpacer(), liveButtonContainer)
			}

			cont.Objects = []fyne.CanvasObject{display}
			cont.Refresh()
		},
	)

	table.SetColumnWidth(0, 120) // Session
	table.SetColumnWidth(1, 120) // Forecast
	table.SetColumnWidth(2, 100) // Date
	table.SetColumnWidth(3, 160) // Time
	table.Resize(fyne.NewSize(520, float32(rowCount*30)))
	return table
}

// ConvertUTCToLocal parses a UTC time string (e.g. "14:30:00Z")
// and a date string (e.g. "2025-04-05") and returns the user's local time.
func ConvertUTCToLocal(dateStr, timeStr string) string {
	// Combine date and time
	datetime := dateStr + "T" + timeStr // e.g. "2025-04-05T14:30:00Z"

	// Parse as RFC3339
	parsed, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return ""
	}

	// Convert to local time zone and format it
	local := parsed.Local()
	return local.Format("3:04 PM MST") // You can adjust the layout
}

func isSessionActive(dateStr, timeStr, sessionType string, now time.Time) bool {
	datetime := dateStr + "T" + timeStr
	parsed, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return false
	}
	start := parsed.Local()

	// Determine session duration and early-show offset
	var duration, early time.Duration

	switch {
	case strings.HasPrefix(sessionType, "Practice"):
		duration = 60 * time.Minute
		early = 15 * time.Minute
	case sessionType == "Qualifying":
		duration = 90 * time.Minute
		early = 30 * time.Minute
	case sessionType == "Sprint":
		duration = 90 * time.Minute
		early = 30 * time.Minute
	case sessionType == "Sprint Qualifying":
		duration = 60 * time.Minute
		early = 15 * time.Minute
	case sessionType == "Race":
		duration = 2 * time.Hour
		early = 60 * time.Minute
	default:
		duration = 0
		early = 0
	}

	visibleStart := start.Add(-early)
	end := start.Add(duration)

	return now.After(visibleStart) && now.Before(end)
}
