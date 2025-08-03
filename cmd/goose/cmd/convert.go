package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/advancedlogic/GoOse/pkg/goose"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [URL]",
	Short: "Extract article content from a URL",
	Long: `Extract clean article content from a web page URL.
This command will fetch the web page, extract the main article content,
and output the cleaned text without advertisements and navigation elements.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// Get flags
		outputFile, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		showMeta, _ := cmd.Flags().GetBool("meta")

		// Extract content
		if err := extractContent(url, outputFile, format, showMeta); err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting content: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Add flags for the convert command
	convertCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	convertCmd.Flags().StringP("format", "f", "text", "Output format: text, json")
	convertCmd.Flags().BoolP("meta", "m", false, "Include metadata (title, author, etc.)")
}

func extractContent(url, outputFile, format string, showMeta bool) error {
	// Create a new GoOse instance
	g := goose.New()

	// Extract article from URL
	article, err := g.ExtractFromURL(url)
	if err != nil {
		return fmt.Errorf("failed to extract article: %w", err)
	}

	var output string

	switch format {
	case "json":
		output = formatJSON(article, showMeta)
	case "text":
		fallthrough
	default:
		output = formatText(article, showMeta)
	}

	// Output to file or stdout
	if outputFile != "" {
		return os.WriteFile(outputFile, []byte(output), 0644)
	} else {
		fmt.Print(output)
		return nil
	}
}

func formatText(article *goose.Article, showMeta bool) string {
	var result string

	if showMeta {
		if article.Title != "" {
			result += fmt.Sprintf("Title: %s\n", article.Title)
		}
		if article.MetaDescription != "" {
			result += fmt.Sprintf("Description: %s\n", article.MetaDescription)
		}
		if article.MetaKeywords != "" {
			result += fmt.Sprintf("Keywords: %s\n", article.MetaKeywords)
		}
		if article.PublishDate != nil {
			result += fmt.Sprintf("Published: %s\n", article.PublishDate.Format("2006-01-02 15:04:05"))
		}
		if article.Domain != "" {
			result += fmt.Sprintf("Domain: %s\n", article.Domain)
		}
		result += "\n"
	}

	// Clean up the extracted text more aggressively
	cleanedText := cleanupExtractedText(article.CleanedText, article.Title)
	result += cleanedText

	return result
}

// cleanupExtractedText performs additional cleanup on extracted text
func cleanupExtractedText(text, title string) string {
	// Split into lines and clean each one
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	seenLines := make(map[string]bool)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines, duplicate lines, and lines that are just the title
		if line == "" {
			// Only add empty line if the last line wasn't empty
			if len(cleanedLines) > 0 && cleanedLines[len(cleanedLines)-1] != "" {
				cleanedLines = append(cleanedLines, "")
			}
		} else if line == title {
			// Skip lines that are just the title repeated
			continue
		} else if !seenLines[line] {
			cleanedLines = append(cleanedLines, line)
			seenLines[line] = true
		}
	}
	
	// Join lines back together
	result := strings.Join(cleanedLines, "\n")
	
	// Remove excessive blank lines
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")
	
	return strings.TrimSpace(result)
}

func formatJSON(article *goose.Article, showMeta bool) string {
	// Simple JSON formatting - you could use encoding/json for more complex cases
	result := "{\n"

	if showMeta {
		if article.Title != "" {
			result += fmt.Sprintf("  \"title\": \"%s\",\n", escapeJSON(article.Title))
		}
		if article.MetaDescription != "" {
			result += fmt.Sprintf("  \"description\": \"%s\",\n", escapeJSON(article.MetaDescription))
		}
		if article.MetaKeywords != "" {
			result += fmt.Sprintf("  \"keywords\": \"%s\",\n", escapeJSON(article.MetaKeywords))
		}
		if article.PublishDate != nil {
			result += fmt.Sprintf("  \"published\": \"%s\",\n", article.PublishDate.Format("2006-01-02T15:04:05Z"))
		}
		if article.Domain != "" {
			result += fmt.Sprintf("  \"domain\": \"%s\",\n", escapeJSON(article.Domain))
		}
	}

	result += fmt.Sprintf("  \"content\": \"%s\"\n", escapeJSON(article.CleanedText))
	result += "}"

	return result
}

func escapeJSON(s string) string {
	// Simple JSON escaping
	s = fmt.Sprintf("%q", s)
	return s[1 : len(s)-1] // Remove surrounding quotes
}
