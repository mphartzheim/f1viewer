package parser

import (
	"encoding/json"

	"github.com/mphartzheim/f1viewer/data"
)

// ParseScheduleResponse parses JSON data into a ScheduleResponse structure.
// Used for endpoints such as ScheduleURL ("https://api.jolpi.ca/ergast/f1/2025.json")
func ParseScheduleResponse(dataBytes []byte) (*data.ScheduleResponse, error) {
	var schedule data.ScheduleResponse
	err := json.Unmarshal(dataBytes, &schedule)
	return &schedule, err
}

// ParseUpcomingResponse parses JSON data into an UpcomingResponse structure.
// Used for endpoints such as UpcomingURL ("https://api.jolpi.ca/ergast/f1/2025/next.json").
func ParseUpcomingResponse(dataBytes []byte) (*data.UpcomingResponse, error) {
	var upcoming data.UpcomingResponse
	err := json.Unmarshal(dataBytes, &upcoming)
	return &upcoming, err
}

// ParseDriverStandingsResponse parses JSON data into a DriverStandingsResponse structure.
// Used for DriversStandingsURL ("https://api.jolpi.ca/ergast/f1/2025/driverstandings.json").
func ParseDriverStandingsResponse(dataBytes []byte) (*data.DriverStandingsResponse, error) {
	var standings data.DriverStandingsResponse
	err := json.Unmarshal(dataBytes, &standings)
	return &standings, err
}

// ParseConstructorStandingsResponse parses JSON data into a ConstructorStandingsResponse structure.
// Used for ConstructorsStandingsURL ("https://api.jolpi.ca/ergast/f1/2025/constructorstandings.json").
func ParseConstructorStandingsResponse(dataBytes []byte) (*data.ConstructorStandingsResponse, error) {
	var standings data.ConstructorStandingsResponse
	err := json.Unmarshal(dataBytes, &standings)
	return &standings, err
}

// ParseRaceResultsResponse parses JSON data into a RaceResultsResponse structure.
// Used for RaceURL ("https://api.jolpi.ca/ergast/f1/2025/2/results.json").
func ParseRaceResultsResponse(dataBytes []byte) (*data.RaceResultsResponse, error) {
	var results data.RaceResultsResponse
	err := json.Unmarshal(dataBytes, &results)
	return &results, err
}

// ParseQualifyingResponse parses JSON data into a QualifyingResponse structure.
// Used for QualifyingURL ("https://api.jolpi.ca/ergast/f1/2025/2/qualifying.json").
func ParseQualifyingResponse(dataBytes []byte) (*data.QualifyingResponse, error) {
	var results data.QualifyingResponse
	err := json.Unmarshal(dataBytes, &results)
	return &results, err
}

// ParseSprintResultsResponse parses JSON data into a SprintResultsResponse structure.
// Used for SprintURL ("https://api.jolpi.ca/ergast/f1/2025/2/sprint.json").
func ParseSprintResultsResponse(dataBytes []byte) (*data.SprintResultsResponse, error) {
	var results data.SprintResultsResponse
	err := json.Unmarshal(dataBytes, &results)
	return &results, err
}
