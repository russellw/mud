package main

// ANSI color codes for text formatting
const (
	// Color codes
	ColorReset     = "\033[0m"
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorUnderline = "\033[4m"
	
	// Foreground colors
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	
	// Bright foreground colors
	ColorBrightBlack   = "\033[90m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"
	
	// Background colors
	ColorBgBlack   = "\033[40m"
	ColorBgRed     = "\033[41m"
	ColorBgGreen   = "\033[42m"
	ColorBgYellow  = "\033[43m"
	ColorBgBlue    = "\033[44m"
	ColorBgMagenta = "\033[45m"
	ColorBgCyan    = "\033[46m"
	ColorBgWhite   = "\033[47m"
)

// Semantic color functions for MUD-specific content
func ColorName(text string) string {
	return ColorBrightCyan + text + ColorReset
}

func ColorRoomName(text string) string {
	return ColorBold + ColorYellow + text + ColorReset
}

func ColorDescription(text string) string {
	return ColorWhite + text + ColorReset
}

func ColorItem(text string) string {
	return ColorBrightGreen + text + ColorReset
}

func ColorMonster(text string) string {
	return ColorBrightRed + text + ColorReset
}

func ColorPlayer(text string) string {
	return ColorBrightBlue + text + ColorReset
}

func ColorExit(text string) string {
	return ColorCyan + text + ColorReset
}

func ColorDamage(text string) string {
	return ColorRed + text + ColorReset
}

func ColorHealing(text string) string {
	return ColorGreen + text + ColorReset
}

func ColorMagic(text string) string {
	return ColorMagenta + text + ColorReset
}

func ColorWarning(text string) string {
	return ColorYellow + text + ColorReset
}

func ColorError(text string) string {
	return ColorBrightRed + text + ColorReset
}

func ColorSuccess(text string) string {
	return ColorBrightGreen + text + ColorReset
}

func ColorInfo(text string) string {
	return ColorBrightCyan + text + ColorReset
}

func ColorEquipment(text string) string {
	return ColorBrightYellow + text + ColorReset
}

// Helper function to colorize text with any color
func Colorize(text, color string) string {
	return color + text + ColorReset
}