package userprefs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
)

// UserPrefs holds the persistent user settings.
type UserPrefs struct {
	Theme        binding.String // e.g., "System", "Dark", "Light", etc.
	Use24hClock  binding.Bool   // Toggle for 24-hour clock display.
	UseLocalTime binding.Bool   // Toggle for local vs. actual/event time.
	HideOnClose  binding.Bool   // Determines whether to hide or quit on close.
	StartHidden  binding.Bool   // Whether the app starts hidden.
}

var prefs *UserPrefs

// Init initializes user preferences by loading from persistent storage,
// or setting defaults if no previous values are stored. It also binds changes
// to auto-persist updates.
func Init() {
	// Get persistent storage using Fyne's Preferences API.
	storage := app.NewWithID("f1viewer").Preferences()

	prefs = &UserPrefs{
		Theme:        binding.NewString(),
		Use24hClock:  binding.NewBool(),
		UseLocalTime: binding.NewBool(),
		HideOnClose:  binding.NewBool(),
		StartHidden:  binding.NewBool(),
	}

	// Helper functions to load or initialize values.
	loadString := func(key string, def string) string {
		val := storage.StringWithFallback(key, def)
		storage.SetString(key, val)
		return val
	}
	loadBool := func(key string, def bool) bool {
		val := storage.BoolWithFallback(key, def)
		storage.SetBool(key, val)
		return val
	}

	// Load saved values or set defaults.
	prefs.Theme.Set(loadString("theme", "System"))
	prefs.Use24hClock.Set(loadBool("use_24h", false))
	prefs.UseLocalTime.Set(loadBool("use_local_time", false))
	prefs.HideOnClose.Set(loadBool("hide_on_close", true))
	prefs.StartHidden.Set(loadBool("start_hidden", false))

	// Auto-persist changes to each setting.
	bindString(storage, "theme", prefs.Theme)
	bindBool(storage, "use_24h", prefs.Use24hClock)
	bindBool(storage, "use_local_time", prefs.UseLocalTime)
	bindBool(storage, "hide_on_close", prefs.HideOnClose)
	bindBool(storage, "start_hidden", prefs.StartHidden)
}

// Get returns the UserPrefs instance for use in your app.
func Get() *UserPrefs {
	return prefs
}

// bindString sets up a listener to persist string changes.
func bindString(storage fyne.Preferences, key string, value binding.String) {
	value.AddListener(binding.NewDataListener(func() {
		if s, err := value.Get(); err == nil {
			storage.SetString(key, s)
		}
	}))
}

// bindBool sets up a listener to persist bool changes.
func bindBool(storage fyne.Preferences, key string, value binding.Bool) {
	value.AddListener(binding.NewDataListener(func() {
		if b, err := value.Get(); err == nil {
			storage.SetBool(key, b)
		}
	}))
}
