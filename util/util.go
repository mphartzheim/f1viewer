package util

import (
	"time"

	"github.com/mphartzheim/f1viewer/userprefs"
)

// FormatTime returns a formatted time string based on the user's 24-hour clock setting.
func FormatTime(t time.Time) string {
	use24, err := userprefs.Get().Use24hClock.Get()
	if err != nil || use24 {
		// If 24h clock is enabled or there's an error, use 24-hour format.
		return t.Format("15:04 MST")
	}
	// Otherwise, use 12-hour format.
	return t.Format("3:04 PM MST")
}
