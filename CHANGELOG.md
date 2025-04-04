All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org).

---

## [v1.0.23] - 2025-04-04

### 🐛 Fixed
- Linux build now correctly installs required OpenGL and X11 development headers (`libgl1-mesa-dev`, `xorg-dev`, and `pkg-config`), resolving Fyne and go-gl compilation issues.

### 🧰 Internal
- Updated `release.yml` to enable CGO and install Linux dependencies before building.
- Improved separation of platform-specific build logic.

---

## [v1.0.22] - 2025-04-04

### 🐛 Fixed
- macOS build now correctly enables CGO and uses the system SDK to support OpenGL dependencies (fixes build failure with `go-gl/gl`).
- Restored macOS compatibility for Fyne-based UI builds.

### 🧰 Internal
- Adjusted `release.yml` to use native builds for macOS instead of cross-compiling.

---

## [v1.0.21] - 2025-04-04

### ✨ Added
- Native Linux packaging: `.deb` and `.rpm` builds are now generated automatically using FPM.
- Checksums are now included for all release archives.

### 🔥 Removed
- AppImage builds. We’re sorry. You didn’t want them. We didn’t need them. They’re gone. (For now.)

### 🧰 Internal
- Updated `release.yml` to support multiplatform packaging and cleaned up legacy logic.

---

## [v1.0.20] - 2025-04-04

### Fixed
- Switched to FUSE method using `fyne package` for Linux AppImage generation
- Installed `fyne` CLI tool in the GitHub Actions environment for packaging

---

## [v1.0.19] - 2025-04-04

### Fixed
- Switched to FUSE method using `fyne package` for Linux AppImage generation
- Ensured `libfuse-dev` and `squashfs-tools` are installed for GitHub Actions builds

---

## [v1.0.18] - 2025-04-04

### Fixed
- Replaced extracted AppImageTool with stable CI-friendly release from TheAssassin/AppImageKit
- Eliminated `exec: : Permission denied` and `.desktop not found` errors
- Now using official prebuilt AppImage binary that works in GitHub Actions

---

## [v1.0.17] - 2025-04-04

### Fixed
- Replaced broken AppImageKit source build with official prebuilt binary
- Extracted `appimagetool-x86_64.AppImage` and used system `mksquashfs`
- Final fix for persistent `/lib/appimagekit/mksquashfs` and CMake errors

---

## [v1.0.16] - 2025-04-04

### Fixed
- Corrected AppImageKit build path by removing incorrect `cd src`
- `cmake` now runs from project root, matching updated AppImageKit structure
- AppImage builds now fully functional with no path, cmake, or FUSE errors

---

## [v1.0.15] - 2025-04-04

### Fixed
- Corrected build failure due to missing `CMakeLists.txt` in AppImageKit root
- Now properly builds `appimagetool` from `AppImageKit/src` with submodules initialized
- Fully functional AppImage generation with no FUSE or squashfs errors remaining

---

## [v1.0.14] - 2025-04-04

### Fixed
- Replaced AppImageTool binary with source-built version to avoid hardcoded FUSE/mksquashfs path issues
- Eliminated `/usr/local/bin/../lib/appimagekit/mksquashfs` error by fully controlling the toolchain
- Now using system squashfs and bypassing all AppImageKit self-contained assumptions

---

## [v1.0.13] - 2025-04-04

### Fixed
- Resolved `mksquashfs` not found error during Linux AppImage build
- Added `squashfs-tools` to GitHub Actions runner to provide `mksquashfs`
- Explicitly exported `MKSQUASHFS` environment variable so `appimagetool` could locate it
- Finalized fully working AppImage pipeline for Linux releases

---

## [v1.0.12] - 2025-04-04

### Fixed
- Installed `squashfs-tools` to support AppImage packaging via `mksquashfs`
- Set `MKSQUASHFS` environment variable so `appimagetool` can find the system squashfs binary
- Resolved missing file error that prevented `.AppImage` creation in GitHub Actions

---

## [v1.0.11] - 2025-04-04

### Fixed
- Corrected AppImage packaging by extracting the proper `appimagetool` binary from the AppImage archive
- Replaced incorrect `AppRun` binary (which caused permission errors) with `squashfs-root/usr/bin/appimagetool`
- Finalized FUSE-free AppImage generation that now builds reliably on GitHub Actions

---

## [v1.0.10] - 2025-04-04

### Fixed
- Replaced AppImage tool installation with a FUSE-free fallback to support GitHub Actions environments
- Extracted `appimagetool` binary from `.AppImage` using `--appimage-extract`, bypassing missing `libfuse.so.2` error
- Ensured Linux AppImage builds work reliably in CI without runtime mounting dependencies

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
