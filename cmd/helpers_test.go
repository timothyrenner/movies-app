package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRuntime(t *testing.T) {
	// Test valid runtime minutes.
	validRuntimeString := "85 min"
	validTruth := 85
	validAnswer, err := ParseRuntime(validRuntimeString)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(validTruth, validAnswer) {
		t.Errorf("Expected %v, got %v", validTruth, validAnswer)
	}

	// Test null runtime minutes.
	nullRuntimeString := "N/A"
	nullTruth := 0
	nullAnswer, err := ParseRuntime(nullRuntimeString)
	if err == nil {
		t.Errorf("Expected error, got nil.")
	}
	if !cmp.Equal(nullTruth, nullAnswer) {
		t.Errorf("Expected %v, got %v", nullTruth, nullAnswer)
	}

	// Test invalid runtime minutes.
	invalidRuntimeString := "31S min"
	invalidTruth := 0
	invalidAnswer, err := ParseRuntime(invalidRuntimeString)
	if err == nil {
		t.Errorf("Expected error, got nil.")
	}
	if !cmp.Equal(invalidTruth, invalidAnswer) {
		t.Errorf("Expected %v, got %v", invalidTruth, invalidAnswer)
	}
}

func TestParseReleased(t *testing.T) {
	// Test "N/A"
	nullReleased := "N/A"
	nullReleasedTruth := "N/A"
	nullReleasedAnswer, err := ParseReleased(nullReleased)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(nullReleasedTruth, nullReleasedAnswer) {
		t.Errorf("Expected %v, got %v", nullReleasedTruth, nullReleasedAnswer)
	}

	// Test a real date.
	released := "08 Sep 2022"
	releasedTruth := "2022-09-08"
	releasedAnswer, err := ParseReleased(released)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(releasedTruth, releasedAnswer) {
		t.Errorf("Expected %v, got %v", releasedTruth, releasedAnswer)
	}
}

func TestSplitOnCommaAndTrim(t *testing.T) {
	toSplit := "Bela Lugosi  ,  Vincent Price,Christopher Lee"
	truth := []string{"Bela Lugosi", "Vincent Price", "Christopher Lee"}
	answer := SplitOnCommaAndTrim(toSplit)
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}
