package src

import (
	"strings"
	"testing"
)

func TestRenderAboutModalSpanish(t *testing.T) {
	m := Model{
		lang:   LangES,
		width:  80,
		height: 24,
	}

	out := m.renderAboutModal()

	expectedSubstrings := []string{
		"bigcrock",
		"Versión:",
		"1.0.0",
		"Creador:",
		"Ars-byte",
		"Repositorio:",
		"github.com/Ars-byte/bigcrock",
		"Licencia:",
		"2026",
		"Hecho con:",
		"Go",
		"Pulsa cualquier tecla para volver",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(out, s) {
			t.Errorf("expected rendered about modal to contain %q, but got:\n%s", s, out)
		}
	}
}

func TestRenderAboutModalEnglish(t *testing.T) {
	m := Model{
		lang:   LangEN,
		width:  80,
		height: 24,
	}

	out := m.renderAboutModal()

	expectedSubstrings := []string{
		"bigcrock",
		"Version:",
		"1.0.0",
		"Creator:",
		"Ars-byte",
		"github.com/Ars-byte/bigcrock",
		"MIT License",
		"2026",
		"Press any key to return",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(out, s) {
			t.Errorf("expected rendered about modal to contain %q, but got:\n%s", s, out)
		}
	}
}
