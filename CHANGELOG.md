# Changelog


- chore(build): update git-cliff integration

- chore(build): additional cliff changes

- chore(build): still testing cliff changes

- chore(build): still testing cliff



- Correct icon in build script
Correct "vv1.0.27" release title
First makefile

- Correct build badge on readme

- Suppress Windows terminal popup by using -H=windowsgui linker flag

- Testing out automated release notes

- makefile now builds changelog and releases code

- Testing automated changelogs



- Build script adjustments
Use Compress-Archive for Windows builds and fix AppImage desktop file location



- Early support for build script
Resolve Windows runner archive step with platform-safe logic



- Upload .deb/.rpm artifacts and enable GitHub release permissions



- Make Go build cross-platform safe by correcting env var handling



- Install OpenGL/X11 build dependencies for Linux release



- Restore macOS build by enabling CGO and system SDK support



- Replaced AppImage with .deb and .rpm builds for Linux



- Install fyne tool in GitHub Actions for AppImage packaging



- Fix AppImage build using Fyne package method



- Fix AppImage build with working CI binary from TheAssassin/AppImageKit



- Fix Linux AppImage build using extracted official appimagetool binary



- Fix AppImageKit cmake step: run cmake from repo root



- Fix AppImageKit build step: use src dir and initialize submodules



- Fix AppImage packaging: build appimagetool from source and use system mksquashfs



- Fix AppImage build: install squashfs-tools and export MKSQUASHFS



- Fix AppImage build: add squashfs-tools and MKSQUASHFS path for packaging



- Fix AppImage tool: use correct extracted binary to avoid exec permission errors



- Fix AppImage build: use extracted appimagetool to avoid FUSE dependency



- Replace fyne package with manual AppImage generation using appimagetool



- Fix AppImage build with pinned Fyne CLI and proper toolchain



- Update AppImage AppID to use GitHub namespace
Use github imports for project



- Fix AppImage upload by splitting release steps per platform



- Fix Windows mkdir, add AppImage renaming and artifact debug output



- Fix AppImage packaging: add appID, release flag, and debug output



- Add AppImage build and fix release permissions



- Fix Fyne CLI install and AppImage support



- Fix release.yml: add Linux deps and simplify jobs

- AppImage support in release actions
Embeds for Windows tray icon



- Initial commit for f1viewer

