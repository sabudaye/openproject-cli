package printer

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

const (
	cBlack   = "\033[30m"
	cRed     = "\033[31m"
	cGreen   = "\033[32m"
	cYellow  = "\033[33m"
	cBlue    = "\033[34m"
	cMagenta = "\033[35m"
	cCyan    = "\033[36m"
	cWhite   = "\033[37m"
	reset    = "\033[0m"
)

// colorsEnabled controls whether the color helpers emit ANSI escape codes.
// It stays on by default so embedded use and tests keep colored output, and
// InitColors turns it off for non-interactive environments: when stdout is
// not a terminal the escape sequences would land as garbage bytes in the
// captured command output.
var colorsEnabled = true

// InitColors decides whether ANSI colors should be emitted, honoring the
// NO_COLOR / CLICOLOR / CLICOLOR_FORCE conventions:
//
//   - NO_COLOR present (any value) -> colors off
//   - CLICOLOR=0                   -> colors off
//   - CLICOLOR_FORCE set to a non-"0" value -> colors on, even without a terminal
//   - otherwise                     -> colors only when stdout is a terminal
func InitColors() {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		colorsEnabled = false
		return
	}
	if force, _ := os.LookupEnv("CLICOLOR_FORCE"); force != "" && force != "0" {
		colorsEnabled = true
		return
	}
	if clicolor, _ := os.LookupEnv("CLICOLOR"); clicolor == "0" {
		colorsEnabled = false
		return
	}

	fd := os.Stdout.Fd()
	colorsEnabled = isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func colorize(code, s string) string {
	if !colorsEnabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", code, s, reset)
}

func Black(s string) string   { return colorize(cBlack, s) }
func Red(s string) string     { return colorize(cRed, s) }
func Green(s string) string   { return colorize(cGreen, s) }
func Yellow(s string) string  { return colorize(cYellow, s) }
func Blue(s string) string    { return colorize(cBlue, s) }
func Magenta(s string) string { return colorize(cMagenta, s) }
func Cyan(s string) string    { return colorize(cCyan, s) }
func White(s string) string   { return colorize(cWhite, s) }
