package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanTitle(t *testing.T) {
	title := " Grizzly 2: Revenge "
	truth := "Grizzly 2 Revenge"
	answer := cleanTitle(title)
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}

	title = "V/H/S 94"
	truth = "VHS 94"
	answer = cleanTitle(title)
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}
