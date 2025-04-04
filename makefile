# Makefile for f1viewer

APPIMAGE_SCRIPT = build/build-appimage.sh

.PHONY: appimage clean changelog release check-cliff

appimage:
	@echo "🚀 Building AppImage..."
	@$(APPIMAGE_SCRIPT)

clean:
	@echo "🧹 Cleaning build directory..."
	@rm -rf build/AppDir build/*.AppImage build/f1viewer*

check-cliff:
	@which git-cliff >/dev/null || (echo "❌ git-cliff not found. Please install it from https://github.com/orhun/git-cliff" && exit 1)

changelog: check-cliff
	@echo "📝 Generating CHANGELOG.md with git-cliff..."
	@git-cliff -o CHANGELOG.md

release: changelog
ifndef VERSION
	$(error VERSION is not set. Usage: make release VERSION=1.0.31)
endif
	@echo "📦 Committing CHANGELOG.md..."
	git add CHANGELOG.md
	git commit -m "docs: update changelog for v$(VERSION)"
	@echo "🏷️  Tagging release v$(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin main
	git push origin v$(VERSION)
