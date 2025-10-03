package main

import (
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("=== Go DocX Library Examples ===")
	log.Println()

	// Example 1: ProcessBytes - Basic placeholder replacement
	log.Println("📄 Example 1: ProcessBytes - Basic Placeholder Replacement")
	log.Println("This demonstrates the original byte-to-byte processing with simple key-value replacements.")
	log.Println()

	// Read template DOCX file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error reading template file: %v", err)
	}

	// Define simple replacements
	replacements := map[string]string{
		"company":     "ACME Corporation",
		"contact":     "John Smith",
		"email":       "john.smith@acme.com",
		"phone":       "(555) 123-4567",
		"address":     "123 Business Ave\nSuite 200\nNew York, NY 10001",
		"website":     "www.acme.com",
		"signature":   "Digital Signature Here",
		"currentYear": "2024",
	}

	// Process the DOCX with standard placeholder replacements
	outputBytes, err := docx.ProcessBytes(templateBytes, replacements)
	if err != nil {
		log.Fatalf("❌ ProcessBytes failed: %v", err)
	}

	// Save the result
	outputPath1 := "basic_processed.docx"
	err = os.WriteFile(outputPath1, outputBytes, 0644)
	if err != nil {
		log.Fatalf("❌ Error writing output file: %v", err)
	}

	log.Printf("✅ Successfully processed with basic replacements and saved to %s", outputPath1)
	log.Printf("   - Static key-value replacements")
	log.Printf("   - Direct text substitution")
	log.Printf("   - No template logic or dynamic content")
	log.Println()

	// Example 2: ProcessTemplate - Template syntax to RTF conversion
	log.Println("🔧 Example 2: ProcessTemplate - Template Syntax to RTF")
	log.Println("This demonstrates converting Go template syntax to RTF format.")
	log.Println()

	// Simple template with basic substitution
	basicTemplate := `Dear {{.CustomerName}},

Thank you for your order #{{.OrderNumber}}.

Your items:
{{range .Items}}- {{.Name}} ({{.Quantity}}x) - ${{printf "%.2f" .Price}}
{{end}}
Subtotal: ${{printf "%.2f" .Subtotal}}
{{if .Tax}}Tax: ${{printf "%.2f" .Tax}}{{end}}
{{if .Discount}}Discount: {{.DiscountPercent}}% ({{printf "%.2f" .DiscountAmount}} off){{end}}
Total: ${{printf "%.2f" .Total}}

{{if .IsRushOrder}}🚀 RUSH ORDER - Express shipping requested{{end}}
Expected delivery: {{.ExpectedDelivery}}

Best regards,
{{.CompanyName}}
{{.CompanyEmail}}`

	// Data for the template
	templateData := map[string]interface{}{
		"CustomerName": "Alice Johnson",
		"OrderNumber":  "ORD-2024-001234",
		"Items": []map[string]interface{}{
			{"Name": "Professional Widget", "Quantity": 2, "Price": 149.99},
			{"Name": "Premium Service", "Quantity": 1, "Price": 299.00},
			{"Name": "Extended Support", "Quantity": 1, "Price": 99.50},
		},
		"Subtotal":         498.48,
		"Tax":              39.88,
		"Discount":         true,
		"DiscountPercent":  10,
		"DiscountAmount":   49.85,
		"Total":            488.51,
		"IsRushOrder":      true,
		"ExpectedDelivery": "Tomorrow by 5 PM",
		"CompanyName":      "Widgets Co.",
		"CompanyEmail":     "orders@widgets.co",
	}

	// Convert template to RTF
	rtfText, err := docx.ProcessTemplate(basicTemplate, templateData)
	if err != nil {
		log.Fatalf("❌ ProcessTemplate failed: %v", err)
	}

	// Save RTF output to file for inspection
	rtfPath := "template_output.rtf"
	err = os.WriteFile(rtfPath, []byte(rtfText), 0644)
	if err != nil {
		log.Fatalf("❌ Error writing RTF file: %v", err)
	}

	log.Printf("✅ Successfully converted template to RTF format")
	log.Printf("   - Complex template with loops, conditions, and formatting")
	log.Printf("   - Dynamic content generation")
	log.Printf("   - RTF output saved to %s", rtfPath)
	log.Printf("   - View the RTF file to see the rich text formatting")
	log.Println()

	// Example 3: ProcessTemplateBytes - Advanced DOCX processing with templates
	log.Println("🚀 Example 3: ProcessTemplateBytes - Advanced DOCX Processing")
	log.Println("This demonstrates integrating template processing with DOCX manipulation.")
	log.Println()

	// Advanced invoice template with conditional logic
	invoiceTemplate := `INVOICE #{{.InvoiceNumber}}

Bill To:
{{.Customer.Name}}
{{.Customer.Address}}
{{.Customer.City}}, {{.Customer.State}} {{.Customer.Zip}}
Email: {{.Customer.Email}}

{{if .Customer.IsBusiness}}Business Tax ID: {{.Customer.TaxID}}{{end}}

Itemized Charges:
{{range .LineItems}}{{printf "%-30s" .Description}}{{printf "%-10s" .Quantity}}{{printf "%10s" (printf "$%.2f" .UnitPrice)}}{{printf "%12s" (printf "$%.2f" .Total)}}\n{{end}}
-------------------------------------------------------
SUBTOTAL:{{printf "%40s" (printf "$%.2f" .Subtotal)}}
{{if .Discount}}DISCOUNT ({{.DiscountPercent}}%):{{printf "%37s" (printf "-$%.2f" .DiscountAmount)}}\n{{end}}{{if .Tax}}TAX:{{printf "%48s" (printf "$%.2f" .Tax)}}\n{{end}}SHIPPING:{{printf "%43s" (printf "$%.2f" .Shipping)}}\nTOTAL DUE:{{printf "%43s" (printf "$%.2f" .GrandTotal)}}

Payment Terms: {{.PaymentTerms}}
Due Date: {{.DueDate}}

{{if .Notes}}Notes: {{.Notes}}

{{end}}Thank you for your business!
{{.Company.Name}}
{{.Company.Address}}
{{.Company.Phone}} | {{.Company.Email}}`

	// Advanced invoice data
	invoiceData := map[string]interface{}{
		"InvoiceNumber": "INV-2024-789",
		"Customer": map[string]interface{}{
			"Name":       "Tech Solutions LLC",
			"Address":    "456 Innovation Drive",
			"City":       "Austin",
			"State":      "TX",
			"Zip":        "78701",
			"Email":      "billing@techsolutions.com",
			"IsBusiness": true,
			"TaxID":      "12-3456789",
		},
		"LineItems": []map[string]interface{}{
			{"Description": "Enterprise Software License", "Quantity": "5", "UnitPrice": 2500.00, "Total": 12500.00},
			{"Description": "Professional Services (100hrs)", "Quantity": "1", "UnitPrice": 150.00, "Total": 15000.00},
			{"Description": "Premium Support Package", "Quantity": "1", "UnitPrice": 2000.00, "Total": 2000.00},
		},
		"Subtotal":        29500.00,
		"Discount":        true,
		"DiscountPercent": 5,
		"DiscountAmount":  1475.00,
		"Tax":             2101.25,
		"Shipping":        100.00,
		"GrandTotal":      30226.25,
		"PaymentTerms":    "Net 30",
		"DueDate":         "2024-02-15",
		"Notes":           "This invoice includes the Q1 enterprise package with custom integrations.",
		"Company": map[string]string{
			"Name":    "Your Business Solutions",
			"Address": "789 Corporate Blvd\nSuite 300\nDallas, TX 75201",
			"Phone":   "(555) 987-6543",
			"Email":   "invoices@yourbusiness.com",
		},
	}

	// Process DOCX with template content
	config := docx.TemplateConfig{
		PlaceholderKey: "key",
		TemplateText:   invoiceTemplate,
		Data:           invoiceData,
	}

	advancedOutput, err := docx.ProcessTemplateBytes(templateBytes, config)
	if err != nil {
		log.Fatalf("❌ ProcessTemplateBytes failed: %v", err)
	}

	// Save advanced template result
	outputPath3 := "advanced_template_processed.docx"
	err = os.WriteFile(outputPath3, advancedOutput, 0644)
	if err != nil {
		log.Fatalf("❌ Error writing advanced output file: %v", err)
	}

	log.Printf("✅ Successfully processed DOCX with advanced template")
	log.Printf("   - Complex business invoice template")
	log.Printf("   - Conditional logic and loops")
	log.Printf("   - Professional formatting")
	log.Printf("   - Output saved to %s", outputPath3)
	log.Println()

	// Summary
	log.Println("📊 Summary of Examples:")
	log.Println(strings.Repeat("=", 25))
	log.Printf("Example 1 (ProcessBytes):     %-50s → Simple key-value replacement", outputPath1)
	log.Printf("Example 2 (ProcessTemplate): %-25s → Template syntax to RTF", rtfPath)
	log.Printf("Example 3 (ProcessTemplateBytes): %-25s → Advanced DOCX processing", outputPath3)
	log.Println()
	log.Println("🎯 Key Differences:")
	log.Println("   • ProcessBytes:       Static replacements, direct file processing")
	log.Println("   • ProcessTemplate:     Template → RTF conversion, no DOCX involved")
	log.Println("   • ProcessTemplateBytes: Template + DOCX processing, most powerful")
	log.Println()
	log.Println("💡 Tips:")
	log.Println("   • Use ProcessBytes for simple placeholder replacement")
	log.Println("   • Use ProcessTemplate for text-only template processing")
	log.Println("   • Use ProcessTemplateBytes for production DOCX generation")
	log.Println("   • All output files are ready to open in Microsoft Word!")
}
