package util

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

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

// ColoredText is a simple widget that displays text in a specific color.
type ColoredText struct {
	widget.BaseWidget
	Text *canvas.Text
}

// NewColoredText creates a new ColoredText widget.
func NewColoredText(textStr string, col color.Color) *ColoredText {
	ct := &ColoredText{
		Text: canvas.NewText(textStr, col),
	}
	ct.ExtendBaseWidget(ct)
	return ct
}

// CreateRenderer implements the fyne.Widget interface.
func (ct *ColoredText) CreateRenderer() fyne.WidgetRenderer {
	// Use a simple renderer that just displays our canvas.Text.
	return widget.NewSimpleRenderer(ct.Text)
}
