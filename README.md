# Go DocX Library

A professional Go library for processing Microsoft Word DOCX files with support for placeholder replacement and template processing.

## Quick Start

```go
import "github.com/izetmolla/go-docx"

// Simple document processing
	doc, err := docx.Open("template.docx")
	if err != nil {
    log.Fatal(err)
}
defer doc.Close()

// Replace placeholders
doc.ReplaceAll(docx.PlaceholderMap{
    "company": "ACME Corp",
    "contact": "John Doe",
})

// Save result
doc.WriteToFile("output.docx")
```

## Three Processing Methods

### 1. ProcessBytes - Simple Replacement
```go
// Load DOCX as bytes
docxBytes, _ := os.ReadFile("template.docx")

// Define replacements
replacements := map[string]string{
    "company": "ACME Corp",
    "contact": "John Doe",
}

// Process and get result
outputBytes, err := docx.ProcessBytes(docxBytes, replacements)
os.WriteFile("output.docx", outputBytes, 0644)
```

### 2. ProcessTemplate - Dynamic Text Generation
```go
// Template with Go syntax
template := `Dear {{.Name}},
Your order #{{.OrderID}} has been {{.Status}}.

{{if .IsVip}}VIP Customer Benefits Applied!{{end}}`

// Data for template
data := map[string]interface{}{
    "Name": "Alice", "OrderID": "12345", 
    "Status": "shipped", "IsVip": true,
}

// Generate plain text output
text, err := docx.ProcessTemplate(template, data)
os.WriteFile("output.txt", []byte(text), 0644)
```

### 3. ProcessTemplateDocx - DOCX Template Processing
```go
// Process DOCX templates containing {{...}} tags directly
docxBytes, err := os.ReadFile("invoice_template.docx")
	if err != nil {
    log.Fatal(err)
}

data := map[string]interface{}{
    "customer": "ABC Corp",
    "items": []map[string]interface{}{
        {"product": "Widget", "quantity": 2, "price": 25.00},
    },
    "total": 50.00,
}

outputBytes, err := docx.ProcessTemplateDocx(docxBytes, data)
os.WriteFile("processed_invoice.docx", outputBytes, 0644)
```

**Note:** Missing data fields preserve original `{{...}}` tags instead of replacing with null/empty values, making it easy to identify incomplete data.

### 4. ProcessTemplateBytes - Advanced DOCX Processing
```go
// Complex template with loops and conditions
template := `Invoice #{{.InvoiceNumber}}
{{range .Items}}
- {{.Product}}: {{.Quantity}}x ${{.Price}}
{{end}}
Total: ${{.Total}}`

// Configure processing
config := docx.TemplateConfig{
    PlaceholderKey: "content",
    TemplateText:   template,
    Data: map[string]interface{}{
        "InvoiceNumber": "INV-001",
        "Items": []map[string]interface{}{
            {"Product": "Widget", "Quantity": 2, "Price": 25.00},
        },
        "Total": 50.00,
    },
}

// Process DOCX with template
outputBytes, err := docx.ProcessTemplateBytes(docxBytes, config)
os.WriteFile("invoice.docx", outputBytes, 0644)
```

## Template Syntax Support

Uses Go's `text/template` package with full syntax support:

```go
{{.Field}}              // Variables
{{if .Condition}}       // Conditions
{{range .Items}}        // Loops
{{printf "%.2f" .Amount}} // Functions
{{/* Comments */}}      // Comments
```

## Examples

Run examples to see different approaches:

```bash
# Quick examples
go run examples/template_simple/main.go

# Comprehensive examples
go run examples/template_basic/main.go

# Method comparison
go run examples/comparison/main.go
```

## When to Use Which Method

| Method | Use For | Output |
|--------|---------|---------|
| **ProcessBytes** | Simple placeholder replacement | DOCX |
| **ProcessTemplate** | Dynamic text generation | Text |
| **ProcessTemplateDocx** | DOCX templates with {{...}} tags | DOCX |
| **ProcessTemplateBytes** | Professional document automation | DOCX |

## Installation

```bash
go get github.com/izetmolla/go-docx
```

## Features

- ✅ Full DOCX document support (headers, footers, images, styles)
- ✅ Template processing with loops and conditions
- ✅ Plain text output for dynamic content generation
- ✅ Byte-to-byte processing for memory efficiency
- ✅ Modern Go 1.24+ with comprehensive error handling

## License

MIT License - see LICENSE file for details.