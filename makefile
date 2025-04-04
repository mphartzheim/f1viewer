# Makefile for f1viewer

APPIMAGE_SCRIPT = build/build-appimage.sh

.PHONY: appimage clean

appimage:
	@echo "🚀 Building AppImage..."
	@$(APPIMAGE_SCRIPT)

clean:
	@echo "🧹 Cleaning build directory..."
	@rm -rf build/AppDir build/*.AppImage build/f1viewer*
