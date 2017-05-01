package goose

import (
	"testing"
)

func TestNormaliseCharset(t *testing.T) {
	characterSet := "SIFT_JIS" // Japanese
	expected := "SHIFT_JIS"
	actual := NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "ANSI" // Western European
	expected = "CP1252"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "ISO-8859-1" // Western European
	expected = "CP1252"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "WINDOWS-1255" // Hebrew
	expected = "ISO-8859-8"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "WINDOWS-1257" // Baltic
	expected = "ISO-8859-13"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "MS949" // Korean
	expected = "UHC"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "KSC5601" // Korean
	expected = "UHC"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "UTF8"
	expected = "UTF-8"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "UTF-8"
	expected = "UTF-8"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "LATIN2"
	expected = "LATIN-2"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}

	characterSet = "WIN1251"
	expected = "CP1251"
	actual = NormaliseCharset(characterSet)
	if expected != actual {
		t.Errorf("Was expecting '%s', got '%s'", expected, actual)
	}
}
