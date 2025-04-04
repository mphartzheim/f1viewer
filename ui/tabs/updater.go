package tabs

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"f1viewer/data"
	"f1viewer/parser"
	"f1viewer/updater"
)

func buildResultsTabContent(title, circuit string, table fyne.CanvasObject) fyne.CanvasObject {
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("%s – %s", title, circuit),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: false},
	)
	return container.NewVScroll(container.NewBorder(header, nil, nil, nil, table))
}

func UpdateRaceResultsTab(season, round string, raceContainer *fyne.Container) {
	ep := updater.Endpoint{
		Name: "Race Results",
		URL:  fmt.Sprintf(data.RaceURL, season, round),
		Parse: func(b []byte) (interface{}, error) {
			return parser.ParseRaceResultsResponse(b)
		},
	}
	result := updater.FetchEndpoint(ep)
	if result.Err != nil {
		fmt.Printf("Race Results endpoint error: %v\n", result.Err)
		raceContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load race results.")}
		raceContainer.Refresh()
		return
	}
	r, ok := result.Data.(*data.RaceResultsResponse)
	if !ok || len(r.MRData.RaceTable.Races) == 0 {
		raceContainer.Objects = []fyne.CanvasObject{widget.NewLabel("No race results available.")}
		raceContainer.Refresh()
		return
	}
	race := r.MRData.RaceTable.Races[0]
	table := CreateRaceResultsTab(r)
	raceContainer.Objects = []fyne.CanvasObject{buildResultsTabContent("Race Results for "+race.RaceName, race.Circuit.CircuitName, table)}
	raceContainer.Refresh()
}

func UpdateQualifyingResultsTab(season, round string, qualifyingContainer *fyne.Container) {
	ep := updater.Endpoint{
		Name: "Qualifying",
		URL:  fmt.Sprintf(data.QualifyingURL, season, round),
		Parse: func(b []byte) (interface{}, error) {
			return parser.ParseQualifyingResponse(b)
		},
	}
	result := updater.FetchEndpoint(ep)
	if result.Err != nil {
		fmt.Printf("Qualifying endpoint error: %v\n", result.Err)
		qualifyingContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load qualifying results.")}
		qualifyingContainer.Refresh()
		return
	}
	r, ok := result.Data.(*data.QualifyingResponse)
	if !ok || len(r.MRData.RaceTable.Races) == 0 {
		qualifyingContainer.Objects = []fyne.CanvasObject{widget.NewLabel("No qualifying results available.")}
		qualifyingContainer.Refresh()
		return
	}
	race := r.MRData.RaceTable.Races[0]
	table := CreateQualifyingResultsTab(r)
	qualifyingContainer.Objects = []fyne.CanvasObject{buildResultsTabContent("Qualifying Results for "+race.RaceName, race.Circuit.CircuitName, table)}
	qualifyingContainer.Refresh()
}

func UpdateSprintResultsTab(season, round string, sprintContainer *fyne.Container) {
	ep := updater.Endpoint{
		Name: "Sprint",
		URL:  fmt.Sprintf(data.SprintURL, season, round),
		Parse: func(b []byte) (interface{}, error) {
			return parser.ParseSprintResultsResponse(b)
		},
	}
	result := updater.FetchEndpoint(ep)
	var raceName, circuitName = "Unknown Race", "Unknown Circuit"
	if result.Err != nil {
		fmt.Printf("Sprint Results endpoint error: %v\n", result.Err)
	} else if r, ok := result.Data.(*data.SprintResultsResponse); ok {
		if len(r.MRData.RaceTable.Races) > 0 {
			race := r.MRData.RaceTable.Races[0]
			raceName = race.RaceName
			circuitName = race.Circuit.CircuitName
			if len(race.SprintResults) > 0 {
				table := CreateSprintResultsTab(r)
				sprintContainer.Objects = []fyne.CanvasObject{buildResultsTabContent("Sprint Results for "+raceName, circuitName, table)}
				sprintContainer.Refresh()
				return
			}
		}
	}
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Sprint Results for %s at %s", raceName, circuitName),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: false},
	)
	body := widget.NewLabel("Not a Sprint Race event.")
	sprintContainer.Objects = []fyne.CanvasObject{container.NewBorder(header, nil, nil, nil, body)}
	sprintContainer.Refresh()
}

func UpdateDriverStandingsTab(season string, driverContainer *fyne.Container) {
	ep := updater.Endpoint{
		Name: "Driver Standings",
		URL:  fmt.Sprintf(data.DriversStandingsURL, season),
		Parse: func(b []byte) (interface{}, error) {
			return parser.ParseDriverStandingsResponse(b)
		},
	}
	result := updater.FetchEndpoint(ep)
	if result.Err != nil {
		fmt.Printf("Driver Standings endpoint error: %v\n", result.Err)
		driverContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load driver standings.")}
		driverContainer.Refresh()
		return
	}
	r, ok := result.Data.(*data.DriverStandingsResponse)
	if !ok || len(r.MRData.StandingsTable.StandingsLists) == 0 {
		driverContainer.Objects = []fyne.CanvasObject{widget.NewLabel("No driver standings available.")}
		driverContainer.Refresh()
		return
	}
	table := CreateDriverStandingsTab(r)
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Driver Standings – %s", season),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: false},
	)
	scroll := container.NewVScroll(container.NewBorder(header, nil, nil, nil, table))
	driverContainer.Objects = []fyne.CanvasObject{scroll}
	driverContainer.Refresh()
}

func UpdateConstructorStandingsTab(season string, constructorContainer *fyne.Container) {
	ep := updater.Endpoint{
		Name: "Constructor Standings",
		URL:  fmt.Sprintf(data.ConstructorsStandingsURL, season),
		Parse: func(b []byte) (interface{}, error) {
			return parser.ParseConstructorStandingsResponse(b)
		},
	}
	result := updater.FetchEndpoint(ep)
	if result.Err != nil {
		fmt.Printf("Constructor Standings endpoint error: %v\n", result.Err)
		constructorContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load constructor standings.")}
		constructorContainer.Refresh()
		return
	}
	r, ok := result.Data.(*data.ConstructorStandingsResponse)
	if !ok || len(r.MRData.StandingsTable.StandingsLists) == 0 {
		constructorContainer.Objects = []fyne.CanvasObject{widget.NewLabel("No constructor standings available.")}
		constructorContainer.Refresh()
		return
	}
	table := CreateConstructorStandingsTab(r)
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Constructor Standings – %s", season),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: false},
	)
	scroll := container.NewVScroll(container.NewBorder(header, nil, nil, nil, table))
	constructorContainer.Objects = []fyne.CanvasObject{scroll}
	constructorContainer.Refresh()
}
