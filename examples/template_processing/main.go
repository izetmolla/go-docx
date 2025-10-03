package main

import (
	"log"
	"os"

	"github.com/izetmolla/go-docx"
)

func main() {
	// Example 1: ProcessTemplate - Convert template syntax to RTF
	// This demonstrates converting Go template syntax to RTF format

	log.Println("=== Example 1: ProcessTemplate ===")
	text := "Dear {{.Name}},\n\nYour order #{{.OrderID}} has been {{.Status}}.\n\nBest regards,\n{{.Company}}"
	data := map[string]string{
		"Name":    "John Doe",
		"OrderID": "12345",
		"Status":  "shipped",
		"Company": "ACME Corp",
	}

	rtfText, err := docx.ProcessTemplate(text, data)
	if err != nil {
		log.Fatalf("ProcessTemplate failed: %v", err)
	}

	log.Printf("RTF Output:\n%s\n", rtfText)

	// Example 2: ProcessTemplateBytes - Process DOCX with template content
	// This demonstrates integrating template processing with DOCX manipulation

	log.Println("=== Example 2: ProcessTemplateBytes ===")

	// Read the template DOCX file
	input, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("Error learning template file: %v", err)
	}

	// Configure template processing
	config := docx.TemplateConfig{
		PlaceholderKey: "key",
		TemplateText:   "Invoice for {{.Customer}}:\n\nItem: {{.Item}}\nQuantity: {{.Quantity}}\nPrice: ${{.Price}}\nTotal: ${{.Total}}\n\nPayment due: {{.DueDate}}",
		Data: map[string]interface{}{
			"Customer": "ABC Corporation",
			"Item":     "Software License",
			"Quantity": "1",
			"Price":    "999.00",
			"Total":    "999.00",
			"DueDate":  "2024-01-31",
		},
	}

	// Process the DOCX with template content
	output, err := docx.ProcessTemplateBytes(input, config)
	if err != nil {
		log.Fatalf("ProcessTemplateBytes failed: %v", err)
	}

	// Write the result to a new file
	outputPath := "template_processed_output.docx"
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	log.Printf("Successfully processed DOCX with template and saved to %s", outputPath)

	// Example 3: Advanced template with conditional logic
	log.Println("=== Example 3: Advanced Template with Conditions ===")

	advancedData := map[string]interface{}{
		"CustomerName": "Jane Smith",
		"IsPremium":    true,
		"Discount":     15,
		"Items": []map[string]interface{}{
			{"Name": "Product A", "Price": 50.00},
			{"Name": "Product B", "Price": 75.50},
		},
	}

	advancedTemplate := `Order Summary for {{.CustomerName}}

{{range .Items}}- {{.Name}}: ${{.Price}}
{{end}}{{if .IsPremium}}Premium Customer Discount: {{.Discount}}%

{{end}}Thank you for your business!

Generated on: {{.OrderDate}}`

	advancedData["OrderDate"] = "2024-01-15"

	advancedRTF, err := docx.ProcessTemplate(advancedTemplate, advancedData)
	if err != nil {
		log.Fatalf("Advanced template processing failed: %v", err)
	}

	log.Printf("Advanced RTF Output:\n%s", advancedRTF)
}
