package cmd

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRuntime(t *testing.T) {
	// Test valid runtime minutes.
	validRuntimeString := "85 min"
	validTruth := sql.NullInt32{
		Int32: 85,
		Valid: true,
	}
	validAnswer := ParseRuntime(validRuntimeString)
	if !cmp.Equal(validTruth, *validAnswer) {
		t.Errorf("Expected %v, got %v", validTruth, *validAnswer)
	}

	// Test null runtime minutes.
	nullRuntimeString := "N/A"
	nullTruth := sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	nullAnswer := ParseRuntime(nullRuntimeString)
	if !cmp.Equal(nullTruth, *nullAnswer) {
		t.Errorf("Expected %v, got %v", nullTruth, *nullAnswer)
	}

	// Test invalid runtime minutes.
	invalidRuntimeString := "31S min"
	invalidTruth := sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	invalidAnswer := ParseRuntime(invalidRuntimeString)
	if !cmp.Equal(invalidTruth, *invalidAnswer) {
		t.Errorf("Expected %v, got %v", invalidTruth, *invalidAnswer)
	}

}
