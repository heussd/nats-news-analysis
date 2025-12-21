package langdetect

import (
	"github.com/pemistahl/lingua-go"
)

var languages = []lingua.Language{
	lingua.English,
	lingua.German,
}

var detector = lingua.NewLanguageDetectorBuilder().
	FromLanguages(languages...).
	WithMinimumRelativeDistance(0.9).
	Build()

func DetectLanguage(text string) lingua.Language {
	language, exists := detector.DetectLanguageOf(text)
	if !exists {
		return lingua.Unknown
	}
	return language
}

type LanguageSubjectPostfix string

const (
	SubjectNewsDe           LanguageSubjectPostfix = "de"
	SubjectNewsEn           LanguageSubjectPostfix = "en"
	SubjectNewsUndetermined LanguageSubjectPostfix = "undetermined" // https://www.loc.gov/standards/iso639-2/faq.html#25
)

func translateLinguaInSubject(language lingua.Language) (subject LanguageSubjectPostfix) {
	switch language {
	case lingua.German:
		subject = SubjectNewsDe
	case lingua.English:
		subject = SubjectNewsEn
	default:
		subject = SubjectNewsUndetermined
	}
	return subject
}

func AssignSubjectPostfixBasedOnLanguage(text string) LanguageSubjectPostfix {
	language := DetectLanguage(text)
	return translateLinguaInSubject(language)
}
