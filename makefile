# Makefile for f1viewer

APPIMAGE_SCRIPT = build/build-appimage.sh

.PHONY: appimage clean changelog release

appimage:
	@echo "🚀 Building AppImage..."
	@$(APPIMAGE_SCRIPT)

clean:
	@echo "🧹 Cleaning build directory..."
	@rm -rf build/AppDir build/*.AppImage build/f1viewer*

changelog:
	@echo "📝 Generating CHANGELOG.md with git-cliff..."
	@git-cliff -o CHANGELOG.md

release: changelog
ifndef VERSION
	$(error VERSION is not set. Usage: make release VERSION=1.0.29)
endif
	@echo "🏷️  Tagging release v$(VERSION)..."
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
