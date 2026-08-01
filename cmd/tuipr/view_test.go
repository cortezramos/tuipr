package main

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 8, "hello..."},
		{"hello", 3, "..."},
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
		}
	}
}

func TestCenterText(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"hello", 10, "  hello"},
		{"hello", 7, " hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hello"},
	}

	for _, tt := range tests {
		result := centerText(tt.input, tt.width)
		if result != tt.expected {
			t.Errorf("centerText(%q, %d) = %q, want %q", tt.input, tt.width, result, tt.expected)
		}
	}
}

func TestGetMergeStatusIcon(t *testing.T) {
	tests := []struct {
		input       string
		expectedStr string
	}{
		{"MERGEABLE", "*"},
		{"CONFLICTING", "W"},
		{"UNKNOWN", "o"},
	}

	for _, tt := range tests {
		icon, _ := getMergeStatusIcon(tt.input)
		if icon != tt.expectedStr {
			t.Errorf("getMergeStatusIcon(%q) icon = %q, want %q", tt.input, icon, tt.expectedStr)
		}
	}
}

func TestGetReviewStatus(t *testing.T) {
	tests := []struct {
		input       string
		shouldExist bool
	}{
		{"APPROVED", true},
		{"CHANGES_REQUESTED", true},
		{"REVIEW_REQUIRED", true},
		{"UNKNOWN", true},
	}

	for _, tt := range tests {
		result := getReviewStatus(tt.input)
		if result == "" && tt.shouldExist {
			t.Errorf("getReviewStatus(%q) returned empty string", tt.input)
		}
	}
}

func TestInfoItem(t *testing.T) {
	tests := []struct {
		name string
		ok   bool
	}{
		{"CI Checks", true},
		{"Conflicts", false},
	}

	for _, tt := range tests {
		result := infoItem(tt.name, tt.ok)
		if result == "" {
			t.Errorf("infoItem(%q, %v) returned empty string", tt.name, tt.ok)
		}
	}
}

func TestSep(t *testing.T) {
	result := sep()
	if result == "" {
		t.Error("sep() should not return empty string")
	}
}

func TestPanelTitle(t *testing.T) {
	result := panelTitle(1, "PRs", true)
	if result == "" {
		t.Error("panelTitle should not return empty string")
	}

	resultInactive := panelTitle(1, "PRs", false)
	if resultInactive == "" {
		t.Error("panelTitle should not return empty string for inactive")
	}
}

func TestPanelIndicator(t *testing.T) {
	result := panelIndicator(1, true)
	if result == "" {
		t.Error("panelIndicator should not return empty string")
	}
}
