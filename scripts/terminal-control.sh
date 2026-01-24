#!/bin/bash
# iTerm/Terminal control script for testing crush themes
# Uses iTerm for test windows to avoid closing the host Terminal session

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SCREENSHOT_DIR="$PROJECT_DIR/screenshots"

# Create screenshot directory
mkdir -p "$SCREENSHOT_DIR"

# Track window IDs we create
CREATED_WINDOWS_FILE="/tmp/crush-test-windows.txt"

# Open a new iTerm window and run a command
iterm_open() {
    local cmd="$1"
    local title="${2:-Crush Test}"

    # Create new iTerm window and get its ID
    osascript <<EOF
tell application "iTerm"
    activate
    delay 0.5

    -- Create a new window
    set newWindow to (create window with default profile)
    delay 0.3

    -- Write command to the session - first cd to project, then run command
    tell current session of newWindow
        write text "cd '$PROJECT_DIR' && $cmd"
    end tell

    -- Try to set window name
    try
        set name of newWindow to "$title"
    end try

    return id of newWindow
end tell
EOF
}

# Send keystrokes to iTerm (control keys)
iterm_send_control() {
    local key="$1"

    osascript \
        -e 'tell application "iTerm" to activate' \
        -e 'delay 0.2' \
        -e 'tell application "System Events"' \
        -e 'tell process "iTerm2"' \
        -e 'keystroke "'"$key"'" using control down' \
        -e 'end tell' \
        -e 'end tell'
}

# Send arrow key to iTerm
iterm_send_arrow() {
    local direction="$1"  # up, down, left, right
    local keycode=""

    case "$direction" in
        up) keycode="126" ;;
        down) keycode="125" ;;
        left) keycode="123" ;;
        right) keycode="124" ;;
    esac

    osascript \
        -e 'tell application "iTerm" to activate' \
        -e 'delay 0.1' \
        -e 'tell application "System Events"' \
        -e 'tell process "iTerm2"' \
        -e "key code $keycode" \
        -e 'end tell' \
        -e 'end tell'
}

# Send Enter key
iterm_send_enter() {
    osascript \
        -e 'tell application "iTerm" to activate' \
        -e 'delay 0.1' \
        -e 'tell application "System Events"' \
        -e 'tell process "iTerm2"' \
        -e 'keystroke return' \
        -e 'end tell' \
        -e 'end tell'
}

# Send Escape key
iterm_send_escape() {
    osascript \
        -e 'tell application "iTerm" to activate' \
        -e 'delay 0.1' \
        -e 'tell application "System Events"' \
        -e 'tell process "iTerm2"' \
        -e 'key code 53' \
        -e 'end tell' \
        -e 'end tell'
}

# Send text to iTerm
iterm_send_text() {
    local text="$1"

    osascript \
        -e 'tell application "iTerm" to activate' \
        -e 'delay 0.1' \
        -e 'tell application "System Events"' \
        -e 'tell process "iTerm2"' \
        -e 'keystroke "'"$text"'"' \
        -e 'end tell' \
        -e 'end tell'
}

# Take a screenshot of iTerm window
iterm_screenshot() {
    local name="${1:-screenshot}"
    local output="$SCREENSHOT_DIR/${name}.png"

    osascript -e 'tell application "iTerm" to activate'
    sleep 0.3

    # Get iTerm window bounds using JSON-style format
    local bounds_str
    bounds_str=$(osascript <<'EOF'
tell application "iTerm"
    set winBounds to bounds of front window
    set x1 to item 1 of winBounds as integer
    set y1 to item 2 of winBounds as integer
    set x2 to item 3 of winBounds as integer
    set y2 to item 4 of winBounds as integer
    return (x1 as string) & " " & (y1 as string) & " " & (x2 as string) & " " & (y2 as string)
end tell
EOF
)

    read -r x1 y1 x2 y2 <<< "$bounds_str"
    local width=$((x2 - x1))
    local height=$((y2 - y1))

    screencapture -x -R"${x1},${y1},${width},${height}" "$output"
    echo "$output"
}

# Close the frontmost iTerm window
iterm_close() {
    osascript \
        -e 'tell application "iTerm"' \
        -e 'close current window' \
        -e 'end tell'
}

# Close all iTerm windows
iterm_close_all() {
    osascript \
        -e 'tell application "iTerm"' \
        -e 'close every window' \
        -e 'end tell'
}

# Quit iTerm entirely
iterm_quit() {
    osascript \
        -e 'tell application "iTerm" to quit'
}

# Resize iTerm window
iterm_resize() {
    local width="${1:-1200}"
    local height="${2:-800}"

    osascript \
        -e 'tell application "iTerm"' \
        -e 'activate' \
        -e "set bounds of front window to {100, 100, $((100 + width)), $((100 + height))}" \
        -e 'end tell'
}

# Wait for app to be ready
iterm_wait() {
    local seconds="${1:-2}"
    sleep "$seconds"
}

# Main command dispatcher
case "$1" in
    open)
        iterm_open "$2" "$3"
        ;;
    control)
        iterm_send_control "$2"
        ;;
    arrow)
        iterm_send_arrow "$2"
        ;;
    enter)
        iterm_send_enter
        ;;
    escape)
        iterm_send_escape
        ;;
    text)
        iterm_send_text "$2"
        ;;
    screenshot)
        iterm_screenshot "$2"
        ;;
    close)
        iterm_close
        ;;
    close-all)
        iterm_close_all
        ;;
    quit)
        iterm_quit
        ;;
    resize)
        iterm_resize "$2" "$3"
        ;;
    wait)
        iterm_wait "$2"
        ;;
    *)
        echo "Usage: $0 {open|control|arrow|enter|escape|text|screenshot|close|close-all|quit|resize|wait}"
        echo ""
        echo "Commands:"
        echo "  open <command> [title]  - Open new iTerm window and run command"
        echo "  control <key>           - Send control+key"
        echo "  arrow <direction>       - Send arrow key (up/down/left/right)"
        echo "  enter                   - Send Enter key"
        echo "  escape                  - Send Escape key"
        echo "  text <string>           - Type text"
        echo "  screenshot [name]       - Take screenshot of iTerm window"
        echo "  close                   - Close front iTerm window"
        echo "  close-all               - Close all iTerm windows"
        echo "  quit                    - Quit iTerm entirely"
        echo "  resize [width] [height] - Resize window (default 1200x800)"
        echo "  wait [seconds]          - Wait for specified seconds"
        exit 1
        ;;
esac
