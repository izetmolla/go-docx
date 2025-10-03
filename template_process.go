package docx

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// ProcessTemplate converts Go template syntax ({{...}}) to plain text.
// This function uses Go's text/template package to parse and execute templates.
//
// Parameters:
//   - templateText: The text containing Go template syntax (e.g., "Hello {{.Name}}")
//   - data: The data object to use for template execution
//
// Returns:
//   - Plain text output with placeholders replaced
//   - An error if template parsing or execution fails
//
// Example:
//
//	text := "Dear {{.Name}},\n\nYour order #{{.OrderID}} has been {{.Status}}.\n\nRegards,\n{{.Company}}"
//	data := map[string]string{
//	    "Name": "John Doe",
//	    "OrderID": "12345",
//	    "Status": "shipped",
//	    "Company": "ACME Corp",
//	}
//	output, err := ProcessTemplate(text, data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// output contains plain text with template variables replaced
func ProcessTemplate(templateText string, data interface{}) (string, error) {
	if templateText == "" {
		return "", fmt.Errorf("template text cannot be empty")
	}

	// Parse the template with text/template
	tmpl, err := template.New("docx_template").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ProcessTemplateBytes combines template processing with DOCX placeholder replacement.
// This function processes a template, then replaces DOCX placeholders with the template output.
//
// Parameters:
//   - input: The DOCX file as bytes
//   - templateConfig: Template configuration containing the template text and data
//
// Returns:
//   - Modified DOCX as bytes with template output
//   - An error if processing fails
//
// Example:
//
//	config := TemplateConfig{
//	    PlaceholderKey: "content",
//	    TemplateText: "Invoice for {{.Customer}}:\nTotal: ${{.Amount}}",
//	    Data: map[string]interface{}{
//	        "Customer": "ABC Corp",
//	        "Amount": "1,234.56",
//	    },
//	}
//	output, err := ProcessTemplateBytes(docxBytes, config)
func ProcessTemplateBytes(input []byte, config TemplateConfig) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("input bytes are empty")
	}

	if config.PlaceholderKey == "" || config.TemplateText == "" {
		return nil, fmt.Errorf("template configuration is incomplete")
	}

	// Process the template
	output, err := ProcessTemplate(config.TemplateText, config.Data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Replace the placeholder in the DOCX with the template output
	replacements := map[string]string{
		config.PlaceholderKey: output,
	}

	return ProcessBytes(input, replacements)
}

// TemplateConfig holds configuration for template processing
type TemplateConfig struct {
	PlaceholderKey string      // The DOCX placeholder key to replace (without delimiters)
	TemplateText   string      // The Go template text with {{...}} syntax
	Data           interface{} // The data object for template execution
}

// ProcessTemplateDocx processes DOCX templates containing Go template syntax ({{...}})
// and replaces them with actual values directly in the Word document.
// This function is similar to ProcessBytes but handles template syntax instead of simple placeholders.
//
// Parameters:
//   - input: The DOCX template file as a byte slice
//   - data: The data object to use for template execution
//
// Returns:
//   - []byte: The processed DOCX document as bytes
//   - error: Any error that occurred during processing
//
// Example:
//
//	docxBytes, err := os.ReadFile("invoice_template.docx")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	data := map[string]interface{}{
//	    "customer": "ABC Corp",
//	    "items": []map[string]interface{}{
//	        {"product": "Widget", "quantity": 2, "price": 25.00},
//	    },
//	    "total": 50.00,
//	}
//
//	outputBytes, err := ProcessTemplateDocx(docxBytes, data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err = os.WriteFile("processed_invoice.docx", outputBytes)
func ProcessTemplateDocx(input []byte, data interface{}) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("input document cannot be empty")
	}

	// Open DOCX from bytes
	doc, err := OpenBytes(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open document from bytes: %w", err)
	}
	defer doc.Close()

	// Find all template placeholders in the document
	placeholders, err := findAllTemplatePlaceholders(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to find template placeholders: %w", err)
	}

	// Process each template placeholder
	replacements := make(PlaceholderMap)
	for placeholder := range placeholders {
		// Execute template for this placeholder
		processedText, err := executeTemplatePlaceholder(placeholder, data)
		if err != nil {
			return nil, fmt.Errorf("failed to process placeholder '%s': %w", placeholder, err)
		}
		replacements[placeholder] = processedText
	}

	// Apply all replacements
	if err := doc.ReplaceAll(replacements); err != nil {
		return nil, fmt.Errorf("failed to replace placeholders: %w", err)
	}

	// Write result to bytes
	var buf bytes.Buffer
	if err := doc.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write document to bytes: %w", err)
	}

	return buf.Bytes(), nil
}

// findAllTemplatePlaceholders searches through the DOCX document for Go template syntax placeholders
// Returns a map of placeholder strings found in the document
func findAllTemplatePlaceholders(doc *Document) (map[string]bool, error) {
	placeholders := make(map[string]bool)

	// Get all placeholders using the existing Document method
	placeholderTexts, err := doc.GetPlaceHoldersList()
	if err != nil {
		return nil, fmt.Errorf("failed to get placeholder list: %w", err)
	}

	// Filter for template syntax ({{...}})
	for _, placeholderText := range placeholderTexts {
		if isTemplatePlaceholder(placeholderText) {
			placeholders[placeholderText] = true
		}
	}

	return placeholders, nil
}

// isTemplatePlaceholder checks if a placeholder contains Go template syntax
func isTemplatePlaceholder(placeholder string) bool {
	return strings.Contains(placeholder, "{{") && strings.Contains(placeholder, "}}")
}

// executeTemplatePlaceholder processes a single template placeholder with the given data
func executeTemplatePlaceholder(templateText string, data interface{}) (string, error) {
	// Parse and execute the template
	tmpl, err := template.New("placeholder").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		// If execution fails due to missing data, return original template text
		// This preserves the {{...}} tags for future identification
		return templateText, nil
	}

	result := buf.String()

	// Check if result indicates missing data (contains placeholder-looking patterns)
	if strings.Contains(result, "<no value>") || strings.Contains(result, "<nil>") ||
		strings.Contains(result, "<invalid Value>") || result == "" {
		// Return original template text to preserve placeholders
		return templateText, nil
	}

	// Check if any {{...}} tags remain in the result (indicating partial processing)
	if strings.Contains(result, "{{") && strings.Contains(result, "}}") {
		// Return original template text to preserve placeholders
		return templateText, nil
	}

	return result, nil
}
