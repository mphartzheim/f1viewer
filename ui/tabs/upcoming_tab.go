package tabs

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
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

	rowCount := len(rows) + 1
	colCount := 3

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

			if id.Row == 0 {
				headers := []string{"Session", "Date", "Time"}
				cont.Objects = []fyne.CanvasObject{widget.NewLabel(headers[id.Col])}
				cont.Refresh()
				return
			}

			row := rows[id.Row-1]
			var display fyne.CanvasObject

			// Attempt to parse and convert to local time
			fullTimeStr := fmt.Sprintf("%sT%s", row.Date, row.Time)
			t, err := time.Parse(time.RFC3339, fullTimeStr)
			localTime := t.Local()
			localDateStr := row.Date // fallback
			localTimeStr := row.Time // fallback
			if err == nil {
				localDateStr = localTime.Format("2006-01-02")
				localTimeStr = localTime.Format("15:04 MST")
			}

			switch id.Col {
			case 0:
				display = widget.NewLabel(row.Session)
			case 1:
				display = widget.NewLabel(localDateStr)
			case 2:
				timeLabel := widget.NewLabel(localTimeStr)
				objects := []fyne.CanvasObject{timeLabel}
				now := time.Now()
				if isSessionActive(row.Date, row.Time, row.Session, now) {
					if u, err := url.Parse(data.F1tvURL); err == nil {
						button := widget.NewButton("🔴 Live", func() {
							_ = fyne.CurrentApp().OpenURL(u)
						})
						button.Importance = widget.HighImportance
						objects = append(objects, button)
					}
				}
				display = container.NewHBox(objects...)
			}

			cont.Objects = []fyne.CanvasObject{display}
			cont.Refresh()
		},
	)

	table.SetColumnWidth(0, 120)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 160)
	table.Resize(fyne.NewSize(400, float32(rowCount*30)))
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
