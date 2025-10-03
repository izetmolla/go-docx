package docx

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

// ProcessTemplate converts Go template syntax ({{...}}) to DOCX format.
// This function uses Go's text/template package to parse and execute templates,
// then formats the output as Rich Text Format (RTF) for Word documents.
//
// Parameters:
//   - templateText: The text containing Go template syntax (e.g., "Hello {{.Name}}")
//   - data: The data object to use for template execution
//
// Returns:
//   - RTF-formatted text that can be inserted into DOCX documents
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
//	rtfText, err := ProcessTemplate(text, data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// rtfText contains RTF-formatted text ready for DOCX insertion
func ProcessTemplate(templateText string, data interface{}) (string, error) {
	if templateText == "" {
		return "", fmt.Errorf("template text cannot be empty")
	}

	// Parse the template with html/template for automatic escaping
	tmpl, err := template.New("docx_template").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Convert the HTML output to RTF format
	rtfOutput, err := convertHTMLToRTF(buf.String())
	if err != nil {
		return "", fmt.Errorf("failed to convert to RTF: %w", err)
	}

	return rtfOutput, nil
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
	rtfOutput, err := ProcessTemplate(config.TemplateText, config.Data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Replace the placeholder in the DOCX with the template output
	replacements := map[string]string{
		config.PlaceholderKey: rtfOutput,
	}

	return ProcessBytes(input, replacements)
}

// TemplateConfig holds configuration for template processing
type TemplateConfig struct {
	PlaceholderKey string      // The DOCX placeholder key to replace (without delimiters)
	TemplateText   string      // The Go template text with {{...}} syntax
	Data           interface{} // The data object for template execution
}

// convertHTMLToRTF converts HTML content to Rich Text Format (RTF) for Word documents.
// This function handles basic HTML formatting and converts it to RTF markup.
func convertHTMLToRTF(htmlContent string) (string, error) {
	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var rtfBuilder strings.Builder

	// Start RTF document
	rtfBuilder.WriteString("{\\rtf1\\ansi\\deff0 {\\fonttbl {\\f0 Times New Roman;}}")

	// Process HTML nodes recursively
	if err := processHTMLNode(doc, &rtfBuilder); err != nil {
		return "", fmt.Errorf("failed to process HTML: %w", err)
	}

	// End RTF document
	rtfBuilder.WriteString("}")

	return rtfBuilder.String(), nil
}

// processHTMLNode recursively processes HTML nodes and converts them to RTF
func processHTMLNode(n *html.Node, rtfBuilder *strings.Builder) error {
	if n.Type == html.TextNode {
		// Handle text content
		text := strings.TrimSpace(n.Data)
		if text != "" {
			// Escape RTF special characters
			escapedText := escapeRTFText(text)
			rtfBuilder.WriteString(escapedText)
		}
	} else if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "br":
			rtfBuilder.WriteString("\\par ")
		case "b", "strong":
			rtfBuilder.WriteString("{\\b ")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := processHTMLNode(c, rtfBuilder); err != nil {
					return err
				}
			}
			rtfBuilder.WriteString("}")
		case "i", "em":
			rtfBuilder.WriteString("{\\i ")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := processHTMLNode(c, rtfBuilder); err != nil {
					return err
				}
			}
			rtfBuilder.WriteString("}")
		case "p":
			rtfBuilder.WriteString("\\par ")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := processHTMLNode(c, rtfBuilder); err != nil {
					return err
				}
			}
			rtfBuilder.WriteString("\\par ")
		case "div":
			// Treat div as paragraph break
			rtfBuilder.WriteString("\\par ")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := processHTMLNode(c, rtfBuilder); err != nil {
					return err
				}
			}
		default:
			// For other elements, process children without special formatting
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := processHTMLNode(c, rtfBuilder); err != nil {
					return err
				}
			}
		}
	} else {
		// Process other node types (like comment nodes)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := processHTMLNode(c, rtfBuilder); err != nil {
				return err
			}
		}
	}

	return nil
}

// escapeRTFText escapes special RTF characters
func escapeRTFText(text string) string {
	// Replace special RTF characters with their escaped versions
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "\\", "\\\\")

	// Handle line breaks - convert \n to RTF line breaks
	text = strings.ReplaceAll(text, "\n", "\\par ")

	return text
}

// ProcessTemplateWithStyles is an enhanced version that supports custom RTF styling
func ProcessTemplateWithStyles(templateText string, data interface{}, styles map[string]string) (string, error) {
	// First process the template normally
	rtfText, err := ProcessTemplate(templateText, data)
	if err != nil {
		return "", err
	}

	// Apply custom styles if provided
	if len(styles) > 0 {
		rtfText = applyCustomStyles(rtfText, styles)
	}

	return rtfText, nil
}

// applyCustomStyles applies custom styles to RTF content
func applyCustomStyles(rtfText string, styles map[string]string) string {
	// Remove the closing } to insert styles
	rtfText = strings.TrimSuffix(rtfText, "}")

	// Add font table if specified
	if fontFamily, ok := styles["font_family"]; ok {
		rtfText += fmt.Sprintf("{\\fonttbl {\\f0 %s;}}", fontFamily)
	}

	// Add color table if specified
	if textColor, ok := styles["text_color"]; ok {
		// Extract RGB values and add to color table
		rtfText += fmt.Sprintf("{\\colortbl %s;}", textColor)
	}

	// Add custom formatting
	if fontSize, ok := styles["font_size"]; ok {
		rtfText = strings.Replace(rtfText, "\\f0", fmt.Sprintf("\\fs%d \\f0", int(parseFloat(fontSize)*2)), 1)
	}

	// Close the RTF document
	rtfText += "}"

	return rtfText
}

// Helper function to parse float from string
func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 12.0 // default font size
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
	for placeholder, _ := range placeholders {
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
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
