#!/bin/bash
set -e

APP="f1viewer"
VERSION=$(git describe --tags --always)
ARCH="x86_64"
OUT_DIR="build"
APPDIR="$OUT_DIR/AppDir"

echo "🔧 Building AppImage for $APP version $VERSION"

# Step 1: Check for dependencies
missing=()
for dep in go appimagetool pkg-config zip; do
  if ! command -v "$dep" &> /dev/null; then
    missing+=("$dep")
  fi
done

if [ ${#missing[@]} -ne 0 ]; then
  echo "❌ Missing dependencies: ${missing[*]}"
  echo "👉 Please install them using dnf:"
  echo "    sudo dnf install ${missing[*]}"
  exit 1
fi

# Step 2: Prepare AppDir
echo "📁 Preparing AppDir structure..."
rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/share/applications"
mkdir -p "$APPDIR/usr/share/icons/hicolor/256x256/apps"

# Step 3: Build the binary
echo "🔨 Building binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o "$APPDIR/usr/bin/$APP"

# Step 4: Add .desktop entry
echo "📝 Creating .desktop file..."
cat > "$APPDIR/usr/share/applications/${APP}.desktop" <<EOF
[Desktop Entry]
Name=$APP
Exec=$APP
Icon=$APP
Type=Application
Categories=Utility;
EOF

# Step 5: Add icon (adjust path if needed)
echo "🖼️  Adding icon..."
cp assets/icon.png "$APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png"

# Step 6: Create AppRun
echo "🚀 Creating AppRun launcher..."
cat > "$APPDIR/AppRun" <<EOF
#!/bin/bash
exec "\$(dirname "\$0")/usr/bin/$APP"
EOF
chmod +x "$APPDIR/AppRun"

# Step 7: Build AppImage
echo "📦 Building AppImage..."
appimagetool "$APPDIR" "$OUT_DIR/${APP}-${VERSION}-${ARCH}.AppImage"

echo "✅ Done! AppImage created at: $OUT_DIR/${APP}-${VERSION}-${ARCH}.AppImage"
