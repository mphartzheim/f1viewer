package main

import (
	_ "embed"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/data"
	"github.com/mphartzheim/f1viewer/parser"
	"github.com/mphartzheim/f1viewer/themes"
	"github.com/mphartzheim/f1viewer/ui/tabs"
	"github.com/mphartzheim/f1viewer/updater"
	"github.com/mphartzheim/f1viewer/userprefs"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
	// Define endpoints.
	endpoints := []updater.Endpoint{
		{
			Name: "Schedule",
			URL:  fmt.Sprintf(data.ScheduleURL, "current"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseScheduleResponse(b)
			},
		},
		{
			Name: "Upcoming",
			URL:  fmt.Sprintf(data.UpcomingURL, "current"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseUpcomingResponse(b)
			},
		},
		{
			Name: "Driver Standings",
			URL:  fmt.Sprintf(data.DriversStandingsURL, "current"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseDriverStandingsResponse(b)
			},
		},
		{
			Name: "Constructor Standings",
			URL:  fmt.Sprintf(data.ConstructorsStandingsURL, "current"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseConstructorStandingsResponse(b)
			},
		},
		{
			Name: "Race Results",
			URL:  fmt.Sprintf(data.RaceURL, "current", "last"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseRaceResultsResponse(b)
			},
		},
		{
			Name: "Qualifying",
			URL:  fmt.Sprintf(data.QualifyingURL, "current", "last"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseQualifyingResponse(b)
			},
		},
		{
			Name: "Sprint Results",
			URL:  fmt.Sprintf(data.SprintURL, "current", "last"),
			Parse: func(b []byte) (interface{}, error) {
				return parser.ParseSprintResultsResponse(b)
			},
		},
	}

	// Initialize user preferences.
	userprefs.Init()
	prefs := userprefs.Get()

	// Map to store last hashes.
	lastHashes := make(map[string]string)

	// Create the Fyne app and main window.
	a := app.NewWithID("ca.jolpi.f1viewer")
	w := a.NewWindow("f1viewer UI")

	// Set the initial theme.
	if themeName, err := prefs.Theme.Get(); err == nil {
		if newTheme, ok := themes.AvailableThemes()[themeName]; ok {
			a.Settings().SetTheme(newTheme)
		}
	}

	// Listen for changes to the theme preference.
	prefs.Theme.AddListener(binding.NewDataListener(func() {
		if themeName, err := prefs.Theme.Get(); err == nil {
			if newTheme, ok := themes.AvailableThemes()[themeName]; ok {
				a.Settings().SetTheme(newTheme)
			}
		}
	}))

	// Tray icon support (only works on desktop platforms)
	if desktopApp, ok := a.(desktop.App); ok {
		icon := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
		maxAttempts := 5
		success := false

		if runtime.GOOS == "windows" {
			for i := 0; i < maxAttempts; i++ {
				func() {
					defer func() { recover() }()
					desktopApp.SetSystemTrayIcon(icon)
					success = true
				}()
				if success {
					break
				}
				fmt.Println("[F1Viewer] Attempt", i+1, "to set system tray icon failed. Retrying...")
				time.Sleep(2 * time.Second)
			}
			if !success {
				fmt.Println("[F1Viewer] Failed to set system tray icon after 5 attempts. Exiting.")
				fyne.CurrentApp().Quit()
				return
			}
		} else {
			desktopApp.SetSystemTrayIcon(icon)
		}

		trayMenu := fyne.NewMenu("F1 Viewer",
			fyne.NewMenuItem("Show Window", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Hide Window", func() {
				w.Hide()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		desktopApp.SetSystemTrayMenu(trayMenu)
		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}

	// Create a top row with a season selector.
	years := generateYears()
	currentYearStr := strconv.Itoa(time.Now().Year())
	seasonSelect := widget.NewSelect(years, func(selected string) {
		fmt.Println("Season selected:", selected)
	})
	seasonSelect.SetSelected(currentYearStr)

	// Create a binding for the countdown timer.
	countdownBinding := binding.NewString()
	countdownBinding.Set("Next: N/A")

	// Top row: Season selector on the left, countdown timer on the right.
	topRow := container.NewHBox(
		widget.NewLabel("Season:"),
		seasonSelect,
		layout.NewSpacer(),
		widget.NewLabelWithData(countdownBinding),
	)

	// Global variables to hold next session info.
	var nextSession time.Time
	var nextRaceName, nextSessionName string

	// Build the outer and inner tab layout.
	outerTabs := container.NewAppTabs()
	scheduleContainer := container.NewStack(widget.NewLabel("Loading schedule..."))
	scheduleTab := container.NewTabItem("Schedule", scheduleContainer)
	outerTabs.Append(scheduleTab)
	upcomingContainer := container.NewStack(widget.NewLabel("Loading upcoming..."))
	upcomingTab := container.NewTabItem("Upcoming", upcomingContainer)
	outerTabs.Append(upcomingTab)
	standingsInnerTabs := container.NewAppTabs()
	driverStandingsContainer := container.NewStack(widget.NewLabel("Loading driver standings..."))
	driverTab := container.NewTabItem("Driver Standings", driverStandingsContainer)
	constructorStandingsContainer := container.NewStack(widget.NewLabel("Loading constructor standings..."))
	constructorTab := container.NewTabItem("Constructor Standings", constructorStandingsContainer)
	standingsInnerTabs.Append(driverTab)
	standingsInnerTabs.Append(constructorTab)
	standingsTab := container.NewTabItem("Standings", standingsInnerTabs)
	outerTabs.Append(standingsTab)
	resultsInnerTabs := container.NewAppTabs()
	raceResultsContainer := container.NewStack(widget.NewLabel("Loading race results..."))
	raceTab := container.NewTabItem("Race", raceResultsContainer)
	qualifyingContainer := container.NewStack(widget.NewLabel("Loading qualifying results..."))
	qualifyingTab := container.NewTabItem("Qualifying", qualifyingContainer)
	sprintResultsContainer := container.NewStack(widget.NewLabel("Loading sprint results..."))
	sprintTab := container.NewTabItem("Sprint", sprintResultsContainer)
	resultsInnerTabs.Append(raceTab)
	resultsInnerTabs.Append(qualifyingTab)
	resultsInnerTabs.Append(sprintTab)
	resultsTab := container.NewTabItem("Results", resultsInnerTabs)
	outerTabs.Append(resultsTab)
	prefsContainer := tabs.CreatePreferencesTab()
	preferencesTab := container.NewTabItem("Preferences", prefsContainer)
	outerTabs.Append(preferencesTab)

	content := container.NewBorder(topRow, nil, nil, nil, outerTabs)
	w.SetContent(content)
	w.Resize(fyne.NewSize(1280, 1024))
	w.SetFixedSize(true)

	// onFlagClicked callback.
	onFlagClicked := func(round string) {
		season := seasonSelect.Selected
		tabs.UpdateRaceResultsTab(season, round, raceResultsContainer)
		tabs.UpdateQualifyingResultsTab(season, round, qualifyingContainer)
		tabs.UpdateSprintResultsTab(season, round, sprintResultsContainer)
		outerTabs.Select(resultsTab)
		resultsInnerTabs.Select(raceTab)
	}

	// updateScheduleTab remains unchanged.
	updateScheduleTab := func() {
		ep := endpoints[0]
		ep.URL = fmt.Sprintf(data.ScheduleURL, seasonSelect.Selected)
		result := updater.FetchEndpoint(ep)
		if result.Err != nil {
			fmt.Printf("Schedule endpoint error: %v\n", result.Err)
			return
		}
		scheduleData, ok := result.Data.(*data.ScheduleResponse)
		if !ok {
			fmt.Println("Failed to cast schedule data")
			return
		}
		table := tabs.CreateScheduleTab(scheduleData, onFlagClicked)
		w.Canvas().Refresh(scheduleContainer)
		scheduleContainer.Objects = []fyne.CanvasObject{table}
		scheduleContainer.Refresh()
	}

	// updateUpcomingTab now determines the next session name.
	updateUpcomingTab := func() {
		ep := endpoints[1]
		epTime := strconv.Itoa(time.Now().Year())
		ep.URL = fmt.Sprintf(data.UpcomingURL, epTime)

		result := updater.FetchEndpoint(ep)
		if result.Err != nil {
			fmt.Printf("Upcoming endpoint error: %v\n", result.Err)
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load upcoming race data."),
			}
			upcomingContainer.Refresh()
			nextSession = time.Time{}
			nextRaceName, nextSessionName = "", ""
			return
		}

		upcomingData, ok := result.Data.(*data.UpcomingResponse)
		if !ok {
			fmt.Println("Failed to cast upcoming data")
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("No upcoming race data found."),
			}
			upcomingContainer.Refresh()
			nextSession = time.Time{}
			nextRaceName, nextSessionName = "", ""
			return
		}

		if len(upcomingData.MRData.RaceTable.Races) > 0 {
			race := upcomingData.MRData.RaceTable.Races[0]
			now := time.Now()

			// Create a helper function to parse a date/time pair.
			parseTimeFunc := func(dateStr, timeStr string) (time.Time, error) {
				var sessionTimeStr string
				if strings.HasSuffix(timeStr, "Z") {
					sessionTimeStr = fmt.Sprintf("%sT%s", dateStr, timeStr)
				} else {
					sessionTimeStr = fmt.Sprintf("%sT%sZ", dateStr, timeStr)
				}
				return time.Parse(time.RFC3339, sessionTimeStr)
			}

			// Build a list of candidate sessions.
			type candidate struct {
				name string
				t    time.Time
			}
			var candidates []candidate

			// Top-level Race.
			if race.Time != "" {
				if t, err := parseTimeFunc(race.Date, race.Time); err == nil {
					candidates = append(candidates, candidate{"Race", t.Local()})
				} else {
					fmt.Printf("Failed to parse Race time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Date, race.Time), err)
				}
			}
			// Practice sessions.
			if race.Practice1.Time != "" {
				if t, err := parseTimeFunc(race.Practice1.Date, race.Practice1.Time); err == nil {
					candidates = append(candidates, candidate{"Practice 1", t.Local()})
				} else {
					fmt.Printf("Failed to parse Practice1 time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Practice1.Date, race.Practice1.Time), err)
				}
			}
			if race.Practice2.Time != "" {
				if t, err := parseTimeFunc(race.Practice2.Date, race.Practice2.Time); err == nil {
					candidates = append(candidates, candidate{"Practice 2", t.Local()})
				} else {
					fmt.Printf("Failed to parse Practice2 time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Practice2.Date, race.Practice2.Time), err)
				}
			}
			if race.Practice3.Time != "" {
				if t, err := parseTimeFunc(race.Practice3.Date, race.Practice3.Time); err == nil {
					candidates = append(candidates, candidate{"Practice 3", t.Local()})
				} else {
					fmt.Printf("Failed to parse Practice3 time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Practice3.Date, race.Practice3.Time), err)
				}
			}
			// Qualifying.
			if race.Qualifying.Time != "" {
				if t, err := parseTimeFunc(race.Qualifying.Date, race.Qualifying.Time); err == nil {
					candidates = append(candidates, candidate{"Qualifying", t.Local()})
				} else {
					fmt.Printf("Failed to parse Qualifying time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Qualifying.Date, race.Qualifying.Time), err)
				}
			}
			// Sprint.
			if race.Sprint.Time != "" {
				if t, err := parseTimeFunc(race.Sprint.Date, race.Sprint.Time); err == nil {
					candidates = append(candidates, candidate{"Sprint", t.Local()})
				} else {
					fmt.Printf("Failed to parse Sprint time: %s, error: %v\n", fmt.Sprintf("%sT%sZ", race.Sprint.Date, race.Sprint.Time), err)
				}
			}

			// Choose the candidate session that is still in the future with the smallest time difference.
			var nextCandidate *candidate
			bestDiff := time.Duration(1<<63 - 1) // a very large duration
			for _, cand := range candidates {
				if cand.t.After(now) {
					diff := cand.t.Sub(now)
					if diff < bestDiff {
						bestDiff = diff
						candCopy := cand
						nextCandidate = &candCopy
					}
				}
			}

			if nextCandidate != nil {
				nextSession = nextCandidate.t
				nextSessionName = nextCandidate.name
				nextRaceName = race.RaceName
			} else {
				nextSession = time.Time{}
				nextRaceName, nextSessionName = race.RaceName, ""
			}

			header := widget.NewLabelWithStyle(
				fmt.Sprintf("Upcoming: %s at %s", race.RaceName, race.Circuit.CircuitName),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: false},
			)
			table := tabs.CreateUpcomingTab(upcomingData)
			wrapped := container.NewBorder(header, nil, nil, nil, table)
			upcomingContainer.Objects = []fyne.CanvasObject{wrapped}
		} else {
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("No upcoming races available."),
			}
			nextSession = time.Time{}
			nextRaceName, nextSessionName = "", ""
		}

		w.Canvas().Refresh(upcomingContainer)
		upcomingContainer.Refresh()
	}

	seasonSelect.OnChanged = func(selected string) {
		outerTabs.Select(scheduleTab)
		updateScheduleTab()
		tabs.UpdateDriverStandingsTab(seasonSelect.Selected, driverStandingsContainer)
		tabs.UpdateConstructorStandingsTab(seasonSelect.Selected, constructorStandingsContainer)
	}

	fmt.Println("Initial load:")
	updater.LoadEndpoints(endpoints, lastHashes)
	updateScheduleTab()

	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		tabs.UpdateRaceResultsTab(seasonSelect.Selected, "last", raceResultsContainer)
	}()
	go func() {
		defer wg.Done()
		tabs.UpdateQualifyingResultsTab(seasonSelect.Selected, "last", qualifyingContainer)
	}()
	go func() {
		defer wg.Done()
		tabs.UpdateSprintResultsTab(seasonSelect.Selected, "last", sprintResultsContainer)
	}()
	go func() {
		defer wg.Done()
		tabs.UpdateDriverStandingsTab(seasonSelect.Selected, driverStandingsContainer)
	}()
	go func() {
		defer wg.Done()
		tabs.UpdateConstructorStandingsTab(seasonSelect.Selected, constructorStandingsContainer)
	}()
	go func() {
		defer wg.Done()
		updateUpcomingTab()
	}()
	wg.Wait()

	// Countdown ticker: update the binding every second.
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			var text string
			if !nextSession.IsZero() {
				remaining := time.Until(nextSession)
				if remaining > 0 {
					totalSeconds := int(remaining.Seconds())
					weeks := totalSeconds / (7 * 24 * 3600)
					days := (totalSeconds % (7 * 24 * 3600)) / (24 * 3600)
					hours := (totalSeconds % (24 * 3600)) / 3600
					minutes := (totalSeconds % 3600) / 60
					seconds := totalSeconds % 60

					if weeks > 0 {
						text = fmt.Sprintf("Next: %s - %s in %dw %dd %02dh %02dm %02ds", nextRaceName, nextSessionName, weeks, days, hours, minutes, seconds)
					} else if days > 0 {
						text = fmt.Sprintf("Next: %s - %s in %dd %02dh %02dm %02ds", nextRaceName, nextSessionName, days, hours, minutes, seconds)
					} else {
						text = fmt.Sprintf("Next: %s - %s in %02dh %02dm %02ds", nextRaceName, nextSessionName, hours, minutes, seconds)
					}
				} else {
					text = "Next session started"
				}
			} else {
				text = "Next: N/A"
			}
			countdownBinding.Set(text)
		}
	}()

	// Refresh the Upcoming tab every 60 seconds.
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			current := outerTabs.Selected()
			if current != nil && current.Text == "Upcoming" {
				fmt.Println("Refreshing Upcoming tab...")
				updateUpcomingTab()
			}
		}
	}()

	w.SetCloseIntercept(func() {
		hide, err := prefs.HideOnClose.Get()
		if err != nil || hide {
			w.Hide()
		} else {
			// If the user prefers to close instead of hiding, close the window and quit the app.
			w.Close()
			a.Quit()
		}
	})

	startHidden, err := prefs.StartHidden.Get()
	if err != nil || !startHidden {
		w.Show()
	}

	a.Run()
}

// Helper: generate years descending from current year down to 1950.
func generateYears() []string {
	currentYear := time.Now().Year()
	years := make([]string, 0, currentYear-1950+1)
	for y := currentYear; y >= 1950; y-- {
		years = append(years, strconv.Itoa(y))
	}
	return years
}
