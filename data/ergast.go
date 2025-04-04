package data

// ----------------------
// Schedule
// ----------------------

// MRDataSchedule is the common outer structure for schedule/upcoming endpoints.
type MRDataSchedule struct {
	XMLNS     string    `json:"xmlns"`
	Series    string    `json:"series"`
	URL       string    `json:"url"`
	Limit     string    `json:"limit"`
	Offset    string    `json:"offset"`
	Total     string    `json:"total"`
	RaceTable RaceTable `json:"RaceTable"`
}

// RaceTable holds the season and its list of races.
type RaceTable struct {
	Season string `json:"season"`
	Races  []Race `json:"Races"`
}

// Race represents a single race in the schedule.
type Race struct {
	Season   string  `json:"season"`
	Round    string  `json:"round"`
	URL      string  `json:"url"`
	RaceName string  `json:"raceName"`
	Circuit  Circuit `json:"Circuit"`
	Date     string  `json:"date"`
	Time     string  `json:"time,omitempty"`
}

// Circuit represents the race circuit details.
type Circuit struct {
	CircuitId   string   `json:"circuitId"`
	URL         string   `json:"url"`
	CircuitName string   `json:"circuitName"`
	Location    Location `json:"Location"`
}

// Location provides geographical details for a circuit.
type Location struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

// ScheduleResponse is used for the full season schedule endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025.json"
type ScheduleResponse struct {
	MRData MRDataSchedule `json:"MRData"`
}

// ----------------------
// Upcoming
// ----------------------

type UpcomingRace struct {
	Season     string      `json:"season"`
	Round      string      `json:"round"`
	URL        string      `json:"url"`
	RaceName   string      `json:"raceName"`
	Circuit    Circuit     `json:"Circuit"`
	Date       string      `json:"date"`
	Time       string      `json:"time,omitempty"`
	Practice1  SessionTime `json:"FirstPractice,omitempty"`
	Practice2  SessionTime `json:"SecondPractice,omitempty"`
	Practice3  SessionTime `json:"ThirdPractice,omitempty"`
	Qualifying SessionTime `json:"Qualifying,omitempty"`
	Sprint     SessionTime `json:"Sprint,omitempty"`
}

type SessionTime struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type UpcomingRaceTable struct {
	Season string         `json:"season"`
	Races  []UpcomingRace `json:"Races"`
}

type MRDataUpcoming struct {
	XMLNS     string            `json:"xmlns"`
	Series    string            `json:"series"`
	URL       string            `json:"url"`
	Limit     string            `json:"limit"`
	Offset    string            `json:"offset"`
	Total     string            `json:"total"`
	RaceTable UpcomingRaceTable `json:"RaceTable"`
}

type UpcomingResponse struct {
	MRData MRDataUpcoming `json:"MRData"`
}

// ----------------------
// Driver Standings
// ----------------------

// MRDataDriverStandings wraps the driver standings response.
type MRDataDriverStandings struct {
	XMLNS          string               `json:"xmlns"`
	Series         string               `json:"series"`
	URL            string               `json:"url"`
	Limit          string               `json:"limit"`
	Offset         string               `json:"offset"`
	Total          string               `json:"total"`
	StandingsTable DriverStandingsTable `json:"StandingsTable"`
}

// DriverStandingsTable holds the standings list for drivers.
type DriverStandingsTable struct {
	Season         string                `json:"season"`
	StandingsLists []DriverStandingsList `json:"StandingsLists"`
}

// DriverStandingsList represents a single set of standings.
type DriverStandingsList struct {
	Season          string           `json:"season"`
	Round           string           `json:"round"`
	DriverStandings []DriverStanding `json:"DriverStandings"`
}

// DriverStanding represents a single driver's standings.
type DriverStanding struct {
	Position     string        `json:"position"`
	Points       string        `json:"points"`
	Wins         string        `json:"wins"`
	Driver       Driver        `json:"Driver"`
	Constructors []Constructor `json:"Constructors"`
}

// Driver represents driver details.
type Driver struct {
	DriverId    string `json:"driverId"`
	URL         string `json:"url"`
	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
}

// DriverStandingsResponse is used for the drivers standings endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025/driverstandings.json"
type DriverStandingsResponse struct {
	MRData MRDataDriverStandings `json:"MRData"`
}

// ----------------------
// Constructor Standings
// ----------------------

// MRDataConstructorStandings wraps the constructor standings response.
type MRDataConstructorStandings struct {
	XMLNS          string                    `json:"xmlns"`
	Series         string                    `json:"series"`
	URL            string                    `json:"url"`
	Limit          string                    `json:"limit"`
	Offset         string                    `json:"offset"`
	Total          string                    `json:"total"`
	StandingsTable ConstructorStandingsTable `json:"StandingsTable"`
}

// ConstructorStandingsTable holds the list for constructors.
type ConstructorStandingsTable struct {
	Season         string                     `json:"season"`
	StandingsLists []ConstructorStandingsList `json:"StandingsLists"`
}

// ConstructorStandingsList represents a standings set.
type ConstructorStandingsList struct {
	Season               string                `json:"season"`
	Round                string                `json:"round"`
	ConstructorStandings []ConstructorStanding `json:"ConstructorStandings"`
}

// ConstructorStanding represents a single constructor's standings.
type ConstructorStanding struct {
	Position    string      `json:"position"`
	Points      string      `json:"points"`
	Wins        string      `json:"wins"`
	Constructor Constructor `json:"Constructor"`
}

// ConstructorStandingsResponse is used for the constructor standings endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025/constructorstandings.json"
type ConstructorStandingsResponse struct {
	MRData MRDataConstructorStandings `json:"MRData"`
}

// Constructor represents a Formula 1 team.
type Constructor struct {
	ConstructorId string `json:"constructorId"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}

// ----------------------
// Race Results
// ----------------------

// MRDataRaceResults wraps the race results response.
type MRDataRaceResults struct {
	XMLNS     string           `json:"xmlns"`
	Series    string           `json:"series"`
	URL       string           `json:"url"`
	Limit     string           `json:"limit"`
	Offset    string           `json:"offset"`
	Total     string           `json:"total"`
	RaceTable RaceResultsTable `json:"RaceTable"`
}

// RaceResultsTable holds the race results.
type RaceResultsTable struct {
	Season string       `json:"season"`
	Round  string       `json:"round"`
	Races  []RaceResult `json:"Races"`
}

// RaceResult represents a specific race and its results.
type RaceResult struct {
	Season   string       `json:"season"`
	Round    string       `json:"round"`
	URL      string       `json:"url"`
	RaceName string       `json:"raceName"`
	Circuit  Circuit      `json:"Circuit"`
	Date     string       `json:"date"`
	Time     string       `json:"time,omitempty"`
	Results  []ResultItem `json:"Results"`
}

// ResultItem represents an individual result in a race.
type ResultItem struct {
	Number       string      `json:"number"`
	Position     string      `json:"position"`
	PositionText string      `json:"positionText"`
	Points       string      `json:"points"`
	Driver       Driver      `json:"Driver"`
	Constructor  Constructor `json:"Constructor"`
	Grid         string      `json:"grid"`
	Laps         string      `json:"laps"`
	Status       string      `json:"status"`
	Time         *RaceTime   `json:"Time,omitempty"`
	FastestLap   *FastestLap `json:"FastestLap,omitempty"`
}

// RaceTime holds the race finishing time.
type RaceTime struct {
	Millis string `json:"millis,omitempty"`
	Time   string `json:"time,omitempty"`
}

// FastestLap holds details about the fastest lap.
type FastestLap struct {
	Rank         string       `json:"rank"`
	Lap          string       `json:"lap"`
	Time         FastLapTime  `json:"Time"`
	AverageSpeed AverageSpeed `json:"AverageSpeed"`
}

// FastLapTime holds the fastest lap time.
type FastLapTime struct {
	Time string `json:"time"`
}

// AverageSpeed holds the average speed details.
type AverageSpeed struct {
	Units string `json:"units"`
	Speed string `json:"speed"`
}

// RaceResultsResponse is used for the race results endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025/2/results.json"
type RaceResultsResponse struct {
	MRData MRDataRaceResults `json:"MRData"`
}

// ----------------------
// Qualifying Results
// ----------------------

// MRDataQualifying wraps the qualifying results response.
type MRDataQualifying struct {
	XMLNS     string          `json:"xmlns"`
	Series    string          `json:"series"`
	URL       string          `json:"url"`
	Limit     string          `json:"limit"`
	Offset    string          `json:"offset"`
	Total     string          `json:"total"`
	RaceTable QualifyingTable `json:"RaceTable"`
}

// QualifyingTable holds qualifying race details.
type QualifyingTable struct {
	Season string           `json:"season"`
	Races  []QualifyingRace `json:"Races"`
}

// QualifyingRace represents a race with qualifying results.
type QualifyingRace struct {
	Season            string             `json:"season"`
	Round             string             `json:"round"`
	URL               string             `json:"url"`
	RaceName          string             `json:"raceName"`
	Circuit           Circuit            `json:"Circuit"`
	QualifyingResults []QualifyingResult `json:"QualifyingResults"`
}

// QualifyingResult represents an individual qualifying result.
type QualifyingResult struct {
	Number      string      `json:"number"`
	Position    string      `json:"position"`
	Driver      Driver      `json:"Driver"`
	Constructor Constructor `json:"Constructor"`
	Q1          string      `json:"Q1,omitempty"`
	Q2          string      `json:"Q2,omitempty"`
	Q3          string      `json:"Q3,omitempty"`
}

// QualifyingResponse is used for the qualifying results endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025/2/qualifying.json"
type QualifyingResponse struct {
	MRData MRDataQualifying `json:"MRData"`
}

// ----------------------
// Sprint Results
// ----------------------

// MRDataSprintResults wraps the sprint results response.
type MRDataSprintResults struct {
	XMLNS     string      `json:"xmlns"`
	Series    string      `json:"series"`
	URL       string      `json:"url"`
	Limit     string      `json:"limit"`
	Offset    string      `json:"offset"`
	Total     string      `json:"total"`
	RaceTable SprintTable `json:"RaceTable"`
}

// SprintTable holds the sprint race details.
type SprintTable struct {
	Season string       `json:"season"`
	Round  string       `json:"round"`
	Races  []SprintRace `json:"Races"`
}

// SprintRace represents a sprint race and its results.
type SprintRace struct {
	Season        string       `json:"season"`
	Round         string       `json:"round"`
	URL           string       `json:"url"`
	RaceName      string       `json:"raceName"`
	Circuit       Circuit      `json:"Circuit"`
	Date          string       `json:"date"`
	Time          string       `json:"time,omitempty"`
	SprintResults []ResultItem `json:"SprintResults"` // <-- corrected key
}

// SprintResultsResponse is used for the sprint results endpoint.
// For example, URL: "https://api.jolpi.ca/ergast/f1/2025/2/sprint.json"
type SprintResultsResponse struct {
	MRData MRDataSprintResults `json:"MRData"`
}
