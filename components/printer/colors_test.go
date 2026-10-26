package printer

import (
	"strings"
	"testing"
)

func TestColorsDisabledWhenNoColorSet(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("CLICOLOR_FORCE", "1")
	InitColors()

	if colorsEnabled {
		t.Error("colors should be disabled when NO_COLOR is set, even with CLICOLOR_FORCE=1")
	}
	if got := Red("x"); strings.ContainsRune(got, '\033') {
		t.Errorf("Red should not emit escape sequences when colors are disabled, got %q", got)
	}
}

func TestColorsDisabledWhenClicolorOff(t *testing.T) {
	t.Setenv("CLICOLOR", "0")
	InitColors()

	if colorsEnabled {
		t.Error("colors should be disabled when CLICOLOR=0")
	}
}

func TestColorsForcedWhenClicolorForce(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")
	InitColors()

	if !colorsEnabled {
		t.Error("colors should be enabled when CLICOLOR_FORCE=1, even without a terminal")
	}
	if got := Cyan("x"); !strings.ContainsRune(got, '\033') {
		t.Errorf("Cyan should emit escape sequences when colors are enabled, got %q", got)
	}
}

// go test captures stdout, so it is never a terminal: without any env
// override InitColors must fall back to the uncolored form.
func TestColorsOffByDefaultWithoutTerminal(t *testing.T) {
	InitColors()

	if colorsEnabled {
		t.Error("colors should be disabled when stdout is not a terminal and no override is set")
	}
}
