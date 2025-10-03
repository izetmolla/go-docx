# Go DocX Library

A professional Go library for processing Microsoft Word DOCX files with support for placeholder replacement.

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

## Processing Methods

### 1. Document API - Full Control
```go
// Open DOCX file
doc, err := docx.Open("template.docx")
if err != nil {
    log.Fatal(err)
}
defer doc.Close()

// Replace individual placeholder
doc.Replace("company", "ACME Corp")

// Replace multiple placeholders
doc.ReplaceAll(docx.PlaceholderMap{
    "company": "ACME Corp",
    "contact": "John Doe",
    "email":   "contact@acme.com",
})

// Save to file
doc.WriteToFile("output.docx")
```

### 2. ProcessBytes - Memory-Efficient Processing
```go
// Load DOCX as bytes
docxBytes, err := os.ReadFile("template.docx")
if err != nil {
    log.Fatal(err)
}

// Define replacements
replacements := map[string]string{
    "company": "ACME Corp",
    "contact": "John Doe",
    "email":   "contact@acme.com",
}

// Process and get result
outputBytes, err := docx.ProcessBytes(docxBytes, replacements)
os.WriteFile("output.docx", outputBytes, 0644)
```

## Examples

Run examples to see different approaches:

```bash
# Simple document processing
go run examples/simple/main.go

# Complex document processing
go run examples/complex/main.go

# Bytes processing
go run examples/bytes_processing/main.go
```

## Installation

```bash
go get github.com/izetmolla/go-docx
```

## Features

- ✅ Full DOCX document support (headers, footers, images, styles)
- ✅ Simple placeholder replacement
- ✅ Memory-efficient byte-to-byte processing
- ✅ Modern Go 1.24+ with comprehensive error handling
- ✅ Cross-platform compatibility

## License

MIT License - see LICENSE file for details.