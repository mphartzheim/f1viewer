All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org).

---

## [v1.0.9] - 2025-04-04

### Fixed
- Replaced `fyne package` with direct AppImage creation using `go build` and `appimagetool`
- Added `.desktop` file and custom AppDir structure to ensure AppImage compatibility
- Fully working `.AppImage` now built and uploaded in CI for Linux releases
- Removed reliance on specific Fyne CLI versions to eliminate silent packaging failures

---

## [v1.0.8] - 2025-04-04

### Fixed
- AppImage generation now works reliably on CI builds
- Installed `appimagetool` in GitHub Actions to support `.AppImage` output
- Pinned Fyne CLI to version `v2.5.5` to avoid CLI argument mismatches and silent packaging failures
- Removed `-verbose` flag which caused `fyne package` to fail on older or mismatched versions

---

## [v1.0.7] - 2025-04-04

### Changed
- Updated AppImage `-appID` to use GitHub project namespace: `com.github.mphartzheim.f1viewer`
- Refactored internal module import paths to use full `github.com/mphartzheim/f1viewer/...` package names
- Updated `go.mod` to reflect new module path for proper Go module resolution

---

## [v1.0.6] - 2025-04-04

### Fixed
- AppImage file was not appearing in GitHub Releases despite successful builds; upload is now handled in a dedicated Linux-only release step
- Release assets are now uploaded per platform (Windows, macOS, Linux) to prevent conditional expression issues
- Improved debug output with `ls dist/` to verify built files before upload

---

## [v1.0.5] - 2025-04-04

### Fixed
- Resolved a workflow failure on Windows caused by `mkdir -p` by switching to a shell-compatible `if` check
- Ensured `.AppImage` is correctly renamed and moved into the `dist/` folder
- Added debug logging to confirm artifact generation during builds
- Improved cross-platform compatibility of the GitHub Actions release process

---

## [v1.0.4] - 2025-04-04

### Fixed
- Corrected `fyne package` usage to ensure `.AppImage` generation
- Added debug logging and fallback handling for AppImage path
- Updated release workflow to prevent CI failures on `mv` step

---

## [v1.0.3] - 2025-04-04

### Fixed
- Added GitHub token permissions to allow release creation
- Prevented `.AppImage` file matching on non-Linux runners

---

## [v1.0.2] - 2025-04-04

### Fixed
- Switched to proper `go install` method for Fyne CLI in Linux build
- Replaced invalid `fyne-cli.zip` download method
- First attempt to generate `.AppImage` for Linux

---

## [v1.0.1] - 2025-04-04

### Added
- Embedded tray icon using `go:embed` for Windows builds

### Changed
- Tray icon loading now uses `fyne.NewStaticResource` instead of disk path
- Internal cleanup to prep release build compatibility

---

## [v1.0.0] - 2025-04-04

### Added
- Full F1 race schedule with local timezone support
- Race, qualifying, and sprint result views
- Driver and constructor standings tabs
- System tray support (Windows, macOS, Linux)
- Tray menu with minimize-to-tray on close
- GitHub Actions setup for cross-platform builds

### Changed
- Schedule tab now highlights the next race dynamically

### Fixed
- Local time conversion bug affecting date display
