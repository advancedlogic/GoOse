package goose

// StopWords is a simple stopwords implementation
type StopWords struct {
	// Simple implementation
}

// NewStopwords creates a new stopwords instance
func NewStopwords() StopWords {
	return StopWords{}
}

// SimpleLanguageDetector detects language (stub implementation)
func (sw StopWords) SimpleLanguageDetector(text string) string {
	return "en" // Default to English
}

// Parser is a simple HTML parser
type Parser struct {
	// Simple implementation
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}
