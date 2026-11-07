package src

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{1073741824 * 10, "10.00 GB"},
		{1073741824 * 150, "150.0 GB"},
	}

	for _, tt := range tests {
		res := FormatBytes(tt.bytes)
		if res != tt.expected {
			t.Errorf("FormatBytes(%d) = %q, expected %q", tt.bytes, res, tt.expected)
		}
	}
}

func TestMakeProgressBar(t *testing.T) {
	bar := MakeProgressBar(0.5, 10)
	if len([]rune(bar)) != 10 {
		t.Errorf("expected bar length 10, got %d", len([]rune(bar)))
	}
	if !strings.HasPrefix(bar, "█████") {
		t.Errorf("expected 5 filled blocks, got %s", bar)
	}

	fullBar := MakeProgressBar(1.0, 8)
	if fullBar != "████████" {
		t.Errorf("expected full bar, got %s", fullBar)
	}

	emptyBar := MakeProgressBar(0.0, 8)
	if emptyBar != "░░░░░░░░" {
		t.Errorf("expected empty bar, got %s", emptyBar)
	}
}

func TestFormatCount(t *testing.T) {
	if FormatCount(500) != "500" {
		t.Errorf("expected '500', got %s", FormatCount(500))
	}
	if FormatCount(1500) != "1,500" {
		t.Errorf("expected '1,500', got %s", FormatCount(1500))
	}
	if FormatCount(1234567) != "1,234,567" {
		t.Errorf("expected '1,234,567', got %s", FormatCount(1234567))
	}
}
