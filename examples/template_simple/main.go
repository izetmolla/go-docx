package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("🚀 Quick Start Examples - Template Processing")
	log.Println(strings.Repeat("=", 45))

	// Read template file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error: %v", err)
	}

	// Example 1: Simple ProcessBytes
	log.Println("📄 Example 1: ProcessBytes (Simple)")
	simpleReplacements := map[string]string{
		"company": "Your Company",
		"contact": "Jane Doe",
		"email":   "jane@yourcompany.com",
	}

	result1, err := docx.ProcessBytes(templateBytes, simpleReplacements)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("simple_output.docx", result1, 0644)
	log.Println("   ✅ Created: simple_output.docx")
	log.Println()

	// Example 2: ProcessTemplate (Text to Plain Text)
	log.Println("🔧 Example 2: ProcessTemplate (Text → Plain Text)")
	template := `Hello {{.Name}},

Your order #{{.OrderID}} for ${{printf "%.2f" .Amount}} has been {{.Status}}.

{{if .ExpressShipping}}🚀 Express shipping - delivery by {{.DeliveryDate}}{{else}}📦 Standard shipping{{end}}

Thank you!`

	data := map[string]interface{}{
		"Name":            "Alice",
		"OrderID":         "12345",
		"Amount":          99.99,
		"Status":          "shipped",
		"ExpressShipping": true,
		"DeliveryDate":    "tomorrow",
	}

	textOutput, err := docx.ProcessTemplate(template, data)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("template_output.txt", []byte(textOutput), 0644)
	log.Println("   ✅ Created: template_output.txt")
	preview := textOutput
	if len(textOutput) > 80 {
		preview = textOutput[:80]
	}
	log.Printf("   📝 Preview: %s...", preview)
	log.Println()

	// Example 3: ProcessTemplateBytes (Advanced)
	log.Println("🚀 Example 3: ProcessTemplateBytes (Advanced)")
	advancedTemplate := `INVOICE #{{.InvoiceNumber}}

{{range .Items}}- {{.Product}}: {{.Quantity}}x ${{.Price|printf "%.2f"}}
{{end}}Total: ${{.Total|printf "%.2f"}}

{{.Customer}} - Thank you!`

	invoiceData := map[string]interface{}{
		"InvoiceNumber": "INV-001",
		"Customer":      "Beta Corp",
		"Items": []map[string]interface{}{
			{"Product": "Widget", "Quantity": 2, "Price": 25.00},
			{"Product": "Gadget", "Quantity": 1, "Price": 50.00},
		},
		"Total": 100.00,
	}

	config := docx.TemplateConfig{
		PlaceholderKey: "key",
		TemplateText:   advancedTemplate,
		Data:           invoiceData,
	}

	result3, err := docx.ProcessTemplateBytes(templateBytes, config)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("advanced_output.docx", result3, 0644)
	log.Println("   ✅ Created: advanced_output.docx")
	log.Println()

	// Summary
	fmt.Printf(`
📊 RESULTS SUMMARY
==================
File                             Method                Output Type
────────────────────────────────────────────────────────────────────
simple_output.docx               ProcessBytes          DOCX
template_output.txt              ProcessTemplate       Text
advanced_output.docx             ProcessTemplateBytes  DOCX

💡 KEY DIFFERENCES:
• ProcessBytes: Static replacements only
• ProcessTemplate: Dynamic text generation → Text  
• ProcessTemplateBytes: Dynamic text generation → DOCX

🎯 WHEN TO USE WHICH:
• ProcessBytes: Simple placeholder replacement
• ProcessTemplate: Text reports, plain text content
• ProcessTemplateBytes: Professional documents

✨ All files ready to open in Microsoft Word!
`)

}
