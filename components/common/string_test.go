package common_test

import (
	"testing"

	"github.com/sabudaye/openproject-cli/components/common"
)

func TestSanitizeLineBreaks(t *testing.T) {
	input := "abc\n"
	expected := "abc"

	result := common.SanitizeLineBreaks(input)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
