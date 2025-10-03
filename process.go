package docx

import (
	"bytes"
	"fmt"
)

// ProcessBytes takes a byte slice representing a DOCX document and a map of
// placeholder replacements. It opens the document from the byte slice,
// performs the replacements, and returns the modified document as a new
// byte slice. This function is ideal for in-memory processing without
// requiring file system operations.
//
// Parameters:
//   - input: The DOCX file as a byte slice
//   - replacements: A map where keys are placeholder names and values are replacement text
//
// Returns:
//   - []byte: The modified DOCX document as bytes
//   - error: Any error that occurred during processing
//
// Example:
//
//	docxBytes, err := os.ReadFile("template.docx")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	replacements := map[string]string{
//	    "company": "ACME Corp",
//	    "contact": "John Doe",
//	}
//
//	outputBytes, err := ProcessBytes(docxBytes, replacements)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err = os.WriteFile("output.docx", outputBytes, 0644)
func ProcessBytes(input []byte, replacements map[string]string) ([]byte, error) {
	doc, err := OpenBytes(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open document from bytes: %w", err)
	}
	defer doc.Close()

	placeholderMap := make(PlaceholderMap)
	for k, v := range replacements {
		placeholderMap[k] = v
	}

	if err := doc.ReplaceAll(placeholderMap); err != nil {
		return nil, fmt.Errorf("failed to replace placeholders: %w", err)
	}

	var buf bytes.Buffer
	if err := doc.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write document to bytes: %w", err)
	}

	return buf.Bytes(), nil
}
