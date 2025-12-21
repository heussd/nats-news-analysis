package langdetect

import (
	"testing"

	"github.com/pemistahl/lingua-go"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected lingua.Language
	}{
		{"This is a test sentence in English.", lingua.English},
		{"Ceci est une phrase de test en français.", lingua.Unknown},
		{"Dies ist ein Testsatz auf Deutsch.", lingua.German},
		{"Esta es una oración de prueba en español.", lingua.Unknown},
		{"これは日本語のテスト文です。", lingua.Unknown},
	}

	for _, test := range tests {
		result := DetectLanguage(test.input)
		if result != test.expected {
			t.Errorf("DetectLanguage(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestAssignSubjectBasedOnLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected LanguageSubjectPostfix
	}{
		{"This is a test sentence in English.", SubjectNewsEn},
		{"Ceci est une phrase de test en français.", SubjectNewsUndetermined},
		{"Dies ist ein Testsatz auf Deutsch.", SubjectNewsDe},
		{"Esta es una oración de prueba en español.", SubjectNewsUndetermined},
		{"これは日本語のテスト文です。", SubjectNewsUndetermined},
	}

	for _, test := range tests {
		result := AssignSubjectPostfixBasedOnLanguage(test.input)
		if result != test.expected {
			t.Errorf("AssignSubjectBasedOnLanguage(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
