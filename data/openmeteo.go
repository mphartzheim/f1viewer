package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mphartzheim/f1viewer/userprefs"
)

type WeatherForecastResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
		Weathercode   []int     `json:"weathercode"`
	} `json:"hourly"`
}

var (
	lastUpcoming     *UpcomingResponse
	lastUpcomingLock sync.RWMutex
)

// SetUpcomingCached stores the latest UpcomingResponse in memory.
func SetUpcomingCached(u *UpcomingResponse) {
	lastUpcomingLock.Lock()
	defer lastUpcomingLock.Unlock()
	lastUpcoming = u
}

// GetUpcomingCached returns the last cached UpcomingResponse.
func GetUpcomingCached() *UpcomingResponse {
	lastUpcomingLock.RLock()
	defer lastUpcomingLock.RUnlock()
	return lastUpcoming
}

func GetForecastForTime(lat, lon float64, t time.Time) (string, error) {
	date := t.Format("2006-01-02")
	url := fmt.Sprintf(WeatherBaseURL, lat, lon, date, date)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var forecast WeatherForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return "", err
	}

	targetHour := t.Format("2006-01-02T15:00")
	for i, ts := range forecast.Hourly.Time {
		if ts == targetHour {
			temp := forecast.Hourly.Temperature2m[i]
			code := forecast.Hourly.Weathercode[i]

			useF, _ := userprefs.Get().UseFahrenheit.Get()
			if useF {
				temp = temp*9/5 + 32
			}
			unit := "°C"
			if useF {
				unit = "°F"
			}

			return fmt.Sprintf("%.0f%s %s", temp, unit, weatherEmoji(code)), nil
		}
	}

	return "No forecast", nil
}

func weatherEmoji(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code >= 1 && code <= 3:
		return "⛅"
	case code >= 45 && code <= 48:
		return "🌫️"
	case code >= 51 && code <= 67:
		return "🌧️"
	case code >= 71 && code <= 77:
		return "🌨️"
	case code >= 80 && code <= 82:
		return "🌦️"
	case code >= 85 && code <= 86:
		return "❄️"
	case code >= 95:
		return "⛈️"
	default:
		return "🌡️"
	}
}
