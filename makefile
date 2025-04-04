# Makefile for f1viewer

APPIMAGE_SCRIPT = build/build-appimage.sh

.PHONY: appimage clean changelog release dry-release check-cliff

appimage: clean
	@echo "🚀 Building AppImage..."
	@$(APPIMAGE_SCRIPT)

clean:
	@echo "🧹 Cleaning build directory..."
	@rm -rf build/AppDir build/*.AppImage build/f1viewer*

check-cliff:
	@which git-cliff >/dev/null || (echo "❌ git-cliff not found. Please install it from https://github.com/orhun/git-cliff" && exit 1)

changelog: check-cliff
ifndef VERSION
	$(error VERSION is not set. Usage: make changelog VERSION=1.0.32)
endif
	@echo "📝 Generating CHANGELOG.md with git-cliff..."
	@git-cliff --tag v$(VERSION) -o CHANGELOG.md -c cliff.toml

release: changelog
ifndef VERSION
	$(error VERSION is not set. Usage: make release VERSION=1.0.32)
endif
	@if git rev-parse "v$(VERSION)" >/dev/null 2>&1; then \
		echo "❌ Tag v$(VERSION) already exists! Use a new version."; \
		exit 1; \
	fi
	@echo "📝 Extracting latest release notes to RELEASENOTES.md..."
	@sed -n "/^## \\[v$(VERSION)\\]/,/^## \\[/p" CHANGELOG.md | sed '$d' > RELEASENOTES.md
	@echo "📦 Committing CHANGELOG.md and RELEASENOTES.md..."
	git add CHANGELOG.md RELEASENOTES.md
	git commit -m "docs: update changelog for v$(VERSION)"
	@echo "🏷️  Tagging release v$(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin main
	git push origin v$(VERSION)

dry-release: changelog
ifndef VERSION
	$(error VERSION is not set. Usage: make dry-release VERSION=1.0.32)
endif
	@echo "🧪 Generating RELEASENOTES.md preview for v$(VERSION)..."
	@sed -n "/^## \\[v$(VERSION)\\]/,/^## \\[/p" CHANGELOG.md | sed '$d' > RELEASENOTES.md
	@echo ""
	@echo "📝 Preview of RELEASENOTES.md:"
	@echo "-----------------------------"
	@cat RELEASENOTES.md
	@echo "-----------------------------"
	@echo "✅ If this looks good, run: make release VERSION=$(VERSION)"
