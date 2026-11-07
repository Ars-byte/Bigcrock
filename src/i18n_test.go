package src

import (
	"testing"
)

func TestLanguageNext(t *testing.T) {
	lang := LangEN
	if lang.Next() != LangES {
		t.Fatalf("expected LangES, got %s", lang.Next())
	}
	if LangES.Next() != LangEN {
		t.Fatalf("expected LangEN, got %s", LangES.Next())
	}
}

func TestDictionariesCompleteness(t *testing.T) {
	languages := []Language{LangEN, LangES}

	for _, lang := range languages {
		d := T(lang)

		if d.AppTitle == "" {
			t.Errorf("[%s] AppTitle is empty", lang)
		}
		if d.AppSubtitle == "" {
			t.Errorf("[%s] AppSubtitle is empty", lang)
		}
		if d.AboutTitle == "" {
			t.Errorf("[%s] AboutTitle is empty", lang)
		}
		if d.AboutVersion == "" {
			t.Errorf("[%s] AboutVersion is empty", lang)
		}
		if d.AboutCreator == "" {
			t.Errorf("[%s] AboutCreator is empty", lang)
		}
		if d.AboutRepo == "" {
			t.Errorf("[%s] AboutRepo is empty", lang)
		}
		if d.AboutLicense == "" {
			t.Errorf("[%s] AboutLicense is empty", lang)
		}
		if d.AboutBuiltWith == "" {
			t.Errorf("[%s] AboutBuiltWith is empty", lang)
		}
		if d.HeaderRank == "" {
			t.Errorf("[%s] HeaderRank is empty", lang)
		}
		if d.HeaderSize == "" {
			t.Errorf("[%s] HeaderSize is empty", lang)
		}
		if d.HeaderUsage == "" {
			t.Errorf("[%s] HeaderUsage is empty", lang)
		}
		if d.HeaderName == "" {
			t.Errorf("[%s] HeaderName is empty", lang)
		}
	}
}
