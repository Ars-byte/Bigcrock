package src

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestNewModel(t *testing.T) {
	m := NewModel(".")
	if m.currDir == "" {
		t.Fatalf("expected non-empty currDir")
	}
	if m.sortMode != SortSizeDesc {
		t.Errorf("expected default sort mode to be SortSizeDesc, got %v", m.sortMode)
	}
	if !m.showHidden {
		t.Errorf("expected showHidden to be true by default")
	}
}

func TestSortItems(t *testing.T) {
	m := Model{
		sortMode: SortSizeDesc,
		items: []Item{
			{Name: "small", Size: 100},
			{Name: "huge", Size: 999999},
			{Name: "medium", Size: 5000},
		},
	}

	m.sortItems()

	if m.items[0].Name != "huge" {
		t.Errorf("expected #1 to be 'huge', got %s", m.items[0].Name)
	}
	if m.items[1].Name != "medium" {
		t.Errorf("expected #2 to be 'medium', got %s", m.items[1].Name)
	}
	if m.items[2].Name != "small" {
		t.Errorf("expected #3 to be 'small', got %s", m.items[2].Name)
	}

	m.sortMode = SortNameAsc
	m.sortItems()

	if m.items[0].Name != "huge" || m.items[1].Name != "medium" || m.items[2].Name != "small" {
		t.Errorf("expected alphabetical ordering")
	}
}

func TestViewResizing(t *testing.T) {
	m := NewModel(".")
	m.items = []Item{
		{Name: ".local", IsDir: true, Size: 1024 * 1024 * 500, SizeKnown: true},
		{Name: "Documents", IsDir: true, Size: 1024 * 1024 * 100, SizeKnown: true},
		{Name: "very_long_file_name_that_could_break_columns_in_narrow_screens.tar.gz", IsDir: false, Size: 1024 * 1024 * 20, SizeKnown: true},
	}

	testSizes := [][2]int{
		{50, 15},
		{70, 20},
		{80, 24},
		{100, 28},
		{120, 32},
		{160, 40},
		{200, 50},
	}

	for _, size := range testSizes {
		w, h := size[0], size[1]
		m.width = w
		m.height = h
		m.adjustScroll()

		rendered := m.View()
		lines := strings.Split(rendered, "\n")

		if len(lines) != h {
			t.Errorf("for size %dx%d, expected %d lines, got %d", w, h, h, len(lines))
		}

		for lineIdx, line := range lines {
			visWidth := ansi.StringWidth(line)
			if visWidth > w {
				t.Errorf("for size %dx%d, line %d visual width %d exceeds window width %d: %q", w, h, lineIdx, visWidth, w, line)
			}
		}
	}
}
