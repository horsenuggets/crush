#!/bin/bash
# Test script for crush themes - takes screenshots of each theme

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SCREENSHOT_DIR="$PROJECT_DIR/screenshots"
TEMP_CONFIG_DIR="/tmp/crush-theme-test"

# All available themes
THEMES=(charmtone dark github hacker light monokai ocean solarized sugarcookie)

# Create directories
mkdir -p "$SCREENSHOT_DIR"
mkdir -p "$TEMP_CONFIG_DIR"

# Clean up function
cleanup() {
    echo "Cleaning up..."
    "$SCRIPT_DIR/terminal-control.sh" close 2>/dev/null || true
    rm -rf "$TEMP_CONFIG_DIR"
}
trap cleanup EXIT

# Function to create config for a theme
create_config() {
    local theme="$1"
    cat > "$TEMP_CONFIG_DIR/crush.json" <<EOF
{
    "options": {
        "tui": {
            "theme": "$theme"
        }
    }
}
EOF
}

# Function to test a single theme
test_theme() {
    local theme="$1"
    echo "Testing theme: $theme"

    # Create config for this theme
    create_config "$theme"

    # Open iTerm with crush using the custom config (need both env vars)
    osascript <<EOF
tell application "iTerm"
    activate
    delay 0.5

    set newWindow to (create window with default profile)
    delay 0.3

    tell current session of newWindow
        write text "cd '$PROJECT_DIR' && CRUSH_GLOBAL_CONFIG='$TEMP_CONFIG_DIR' CRUSH_GLOBAL_DATA='$TEMP_CONFIG_DIR' ./crush-test"
    end tell
end tell
EOF

    # Wait for crush to start
    sleep 4

    # Resize window for consistent screenshots
    "$SCRIPT_DIR/terminal-control.sh" resize 1200 800
    sleep 0.5

    # Take screenshot
    "$SCRIPT_DIR/terminal-control.sh" screenshot "$theme"
    echo "Screenshot saved: $SCREENSHOT_DIR/$theme.png"

    # Close the window
    "$SCRIPT_DIR/terminal-control.sh" close
    sleep 1
}

# Main execution
echo "Starting theme testing..."
echo "Screenshots will be saved to: $SCREENSHOT_DIR"
echo ""

# Test each theme
for theme in "${THEMES[@]}"; do
    test_theme "$theme"
done

echo ""
echo "All themes tested. Screenshots saved to $SCREENSHOT_DIR"
ls -la "$SCREENSHOT_DIR"
