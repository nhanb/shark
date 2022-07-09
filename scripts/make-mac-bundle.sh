#!/usr/bin/env sh
set -euf

echo "Current version: $TAG"

MAC_EXECUTABLE="$1"
MAC_ICON=icon.icns
MAC_BUNDLE_DIR=dist/Shark.app

mkdir -p $MAC_BUNDLE_DIR/Contents/Resources
mkdir -p $MAC_BUNDLE_DIR/Contents/MacOS

cp "$MAC_ICON" $MAC_BUNDLE_DIR/Contents/Resources/Shark.icns
cp "$MAC_EXECUTABLE" $MAC_BUNDLE_DIR/Contents/MacOS/Shark
chmod +x $MAC_BUNDLE_DIR/Contents/MacOS/Shark

cat <<EOT >> $MAC_BUNDLE_DIR/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8" standalone="no"?><plist version="1.0">
  <dict>
    <key>CFBundleExecutable</key>
    <string>Shark</string>
    <key>CFBundleGetInfoString</key>
    <string>Shark $TAG</string>
    <key>CFBundleVersion</key>
    <string>0.2</string>
    <key>CFBundleShortVersionString</key>
    <string>0.2</string>
    <key>CFBundleIconFile</key>
    <string>Shark</string>
    <key>CFBundleIdentifier</key>
    <string>com.imnhan.shark</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOT

rm "$MAC_EXECUTABLE"
