package cmd

var (
	Gray    = "\033[90m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Reset   = "\033[0m"
)

var LevelColors = map[int]string{
	0:  Blue,
	10: Gray,    // Trace
	20: Yellow,  // Debug
	30: Cyan,    // Info
	40: Magenta, // Warning
	50: Red,     // Error
	60: Red,     // Fatal
}
