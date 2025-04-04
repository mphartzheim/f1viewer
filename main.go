package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
	"f1viewer/parser"
	"f1viewer/ui/tabs"
	"f1viewer/updater"
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

	// Map to store last hashes.
	lastHashes := make(map[string]string)

	// Create the Fyne app and main window.
	a := app.NewWithID("ca.jolpi.f1viewer")
	w := a.NewWindow("f1viewer UI")

	// Tray icon support (only works on desktop platforms)
	if desktopApp, ok := a.(desktop.App); ok {
		// Load embedded tray icon
		icon := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
		desktopApp.SetSystemTrayIcon(icon)

		// Tray menu with show/hide and quit
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

		// Intercept close to hide instead of quit
		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}

	// Create a top row with a season selector.
	years := generateYears()
	currentYearStr := strconv.Itoa(time.Now().Year())
	seasonSelect := widget.NewSelect(years, func(selected string) {
		fmt.Println("Season selected:", selected)
		// Reload endpoints with the selected season if needed.
	})
	seasonSelect.SetSelected(currentYearStr)
	topRow := container.NewHBox(widget.NewLabel("Season:"), seasonSelect)

	// Build the outer and inner tab layout.
	outerTabs := container.NewAppTabs()

	// Schedule Tab.
	scheduleContainer := container.NewStack(widget.NewLabel("Loading schedule..."))
	scheduleTab := container.NewTabItem("Schedule", scheduleContainer)
	outerTabs.Append(scheduleTab)

	// Upcoming Tab.
	// Now the Upcoming tab loads a table from our new CreateUpcomingTab function.
	upcomingContainer := container.NewStack(widget.NewLabel("Loading upcoming..."))
	upcomingTab := container.NewTabItem("Upcoming", upcomingContainer)
	outerTabs.Append(upcomingTab)

	// Standings Tab with inner tabs.
	standingsInnerTabs := container.NewAppTabs()
	driverStandingsContainer := container.NewStack(widget.NewLabel("Loading driver standings..."))
	driverTab := container.NewTabItem("Driver Standings", driverStandingsContainer)
	constructorStandingsContainer := container.NewStack(widget.NewLabel("Loading constructor standings..."))
	constructorTab := container.NewTabItem("Constructor Standings", constructorStandingsContainer)
	standingsInnerTabs.Append(driverTab)
	standingsInnerTabs.Append(constructorTab)
	standingsTab := container.NewTabItem("Standings", standingsInnerTabs)
	outerTabs.Append(standingsTab)

	// Results Tab with inner tabs.
	resultsInnerTabs := container.NewAppTabs()

	// Race Results container.
	raceResultsContainer := container.NewStack(widget.NewLabel("Loading race results..."))
	raceTab := container.NewTabItem("Race", raceResultsContainer)

	// Qualifying Results container.
	qualifyingContainer := container.NewStack(widget.NewLabel("Loading qualifying results..."))
	qualifyingTab := container.NewTabItem("Qualifying", qualifyingContainer)

	// Sprint Results container.
	sprintResultsContainer := container.NewStack(widget.NewLabel("Loading sprint results..."))
	sprintTab := container.NewTabItem("Sprint", sprintResultsContainer)

	resultsInnerTabs.Append(raceTab)
	resultsInnerTabs.Append(qualifyingTab)
	resultsInnerTabs.Append(sprintTab)

	resultsTab := container.NewTabItem("Results", resultsInnerTabs)
	outerTabs.Append(resultsTab)

	// Combine topRow and outerTabs in a border layout.
	content := container.NewBorder(topRow, nil, nil, nil, outerTabs)
	w.SetContent(content)
	w.Resize(fyne.NewSize(1280, 1024))
	w.SetFixedSize(true)
	w.Show()

	// onFlagClicked callback: when flag is clicked in the Schedule tab.
	onFlagClicked := func(round string) {
		season := seasonSelect.Selected

		// --- Load Results ---
		tabs.UpdateRaceResultsTab(season, round, raceResultsContainer)
		tabs.UpdateQualifyingResultsTab(season, round, qualifyingContainer)
		tabs.UpdateSprintResultsTab(season, round, sprintResultsContainer)

		// Show outer tab + Race as default selection.
		outerTabs.Select(resultsTab)
		resultsInnerTabs.Select(raceTab)
	}

	// updateScheduleTab fetches schedule data and updates scheduleContainer.
	// It now passes the onFlagClicked callback to createScheduleTab.
	updateScheduleTab := func() {
		ep := endpoints[0]
		// Optionally, update the endpoint URL based on seasonSelect.Selected.
		ep.URL = fmt.Sprintf(data.ScheduleURL, seasonSelect.Selected)
		result := updater.FetchEndpoint(ep) // "Schedule" endpoint.
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

	// updateUpcomingTab fetches upcoming data and updates upcomingContainer.
	updateUpcomingTab := func() {
		ep := endpoints[1] // Upcoming endpoint
		epTime := strconv.Itoa(time.Now().Year())
		ep.URL = fmt.Sprintf(data.UpcomingURL, epTime)

		result := updater.FetchEndpoint(ep)
		if result.Err != nil {
			fmt.Printf("Upcoming endpoint error: %v\n", result.Err)
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load upcoming race data."),
			}
			upcomingContainer.Refresh()
			return
		}

		upcomingData, ok := result.Data.(*data.UpcomingResponse)
		if !ok {
			fmt.Println("Failed to cast upcoming data")
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("No upcoming race data found."),
			}
			upcomingContainer.Refresh()
			return
		}

		if len(upcomingData.MRData.RaceTable.Races) > 0 {
			race := upcomingData.MRData.RaceTable.Races[0]
			raceName := race.RaceName
			circuitName := race.Circuit.CircuitName

			header := widget.NewLabelWithStyle(
				fmt.Sprintf("Upcoming: %s at %s", raceName, circuitName),
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
		}

		w.Canvas().Refresh(upcomingContainer)
		upcomingContainer.Refresh()
	}

	// Hook season selector: update the schedule tab when a new season is selected.
	seasonSelect.OnChanged = func(selected string) {
		outerTabs.Select(scheduleTab)
		updateScheduleTab()
		tabs.UpdateDriverStandingsTab(seasonSelect.Selected, driverStandingsContainer)
		tabs.UpdateConstructorStandingsTab(seasonSelect.Selected, constructorStandingsContainer)
	}

	// Hook Upcoming tab: only load when the Upcoming tab is selected.
	outerTabs.OnSelected = func(selected *container.TabItem) {
		if selected.Text == "Upcoming" {
			updateUpcomingTab()
		}
	}

	// Initial load.
	fmt.Println("Initial load:")
	updater.LoadEndpoints(endpoints, lastHashes)
	updateScheduleTab()

	var wg sync.WaitGroup
	wg.Add(5) // We're launching 5 goroutines

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

	wg.Wait()

	// Update the Upcoming tab every 60 seconds.
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
