package main

import (
	"log"
	"os"

	"github.com/izetmolla/go-docx"
)

func main() {
	// Load the template DOCX file into memory
	templateBytes, err := os.ReadFile("../simple/template.docx")
	if err != nil {
		log.Fatalf("Failed to read template: %v", err)
	}

	// Define replacements (keys without delimiters)
	replacements := map[string]string{
		"key":                         "Professional Golang",
		"key-with-dash":               "go-docx",
		"key-with-dashes":             "byte-processing",
		"key with space":              "Modern API",
		"key_with_underscore":         "Clean_Code",
		"multiline":                   "First Line\nSecond Line\nThird Line",
		"key.with.dots":               "v1.0.0",
		"mixed-key.separator_styles#": "Professional",
		"yet-another_placeholder":     "Refactored",
		"foo":                         "bar",
		"ampersand":                   "Fast & Reliable",
		"newlinetester":               "Line 1\nLine 2",
	}

	// Process the DOCX entirely in memory (byte-to-byte)
	outputBytes, err := docx.ProcessBytes(templateBytes, replacements)
	if err != nil {
		log.Fatalf("Failed to process DOCX: %v", err)
	}

	// Write the result to a new file
	err = os.WriteFile("output_bytes.docx", outputBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	log.Printf("Successfully processed DOCX using byte-to-byte processing")
	log.Printf("Input size: %d bytes", len(templateBytes))
	log.Printf("Output size: %d bytes", len(outputBytes))
	log.Printf("Output saved to: output_bytes.docx")
}
