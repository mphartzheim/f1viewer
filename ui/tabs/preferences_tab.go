package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1viewer/themes"
	"github.com/mphartzheim/f1viewer/userprefs"
)

// CreatePreferencesTab builds the Preferences tab UI. It allows users to
// update settings that are persisted across app loads. Options include:
//   - Theme selection from a custom registry.
//   - Toggle for a 24-hour clock display.
//   - Toggle for displaying local vs. event time.
//   - Toggle for whether the window hides on close.
//   - Toggle for whether the app starts hidden.
func CreatePreferencesTab() fyne.CanvasObject {
	prefs := userprefs.Get()

	// Create theme selection dropdown using our themes registry.
	themeOptions := themes.SortedThemeList()
	themeSelect := widget.NewSelect(themeOptions, func(selected string) {
		prefs.Theme.Set(selected)
	})
	// Initialize the dropdown with the current theme setting.
	if currentTheme, err := prefs.Theme.Get(); err == nil {
		themeSelect.SetSelected(currentTheme)
	}
	// Wrap the select widget in an HBox to ensure it only takes up as much space as needed.
	themeContainer := container.NewHBox(themeSelect)

	// Create toggle checkboxes with bindings for persistence.
	clockCheck := widget.NewCheckWithData("Use 24-Hour Clock", prefs.Use24hClock)
	localTimeCheck := widget.NewCheckWithData("Display Local Time", prefs.UseLocalTime)
	hideOnCloseCheck := widget.NewCheckWithData("Hide on Close", prefs.HideOnClose)
	startHiddenCheck := widget.NewCheckWithData("Start Hidden", prefs.StartHidden)

	// Build a form to group the preference settings.
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: themeContainer},
			{Text: "Clock", Widget: clockCheck},
			{Text: "Display Time", Widget: localTimeCheck},
			{Text: "Close Behavior", Widget: hideOnCloseCheck},
			{Text: "Start Behavior", Widget: startHiddenCheck},
		},
	}

	// Add a header for the tab.
	header := widget.NewLabelWithStyle("Preferences", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	return container.NewVBox(header, form)
}
