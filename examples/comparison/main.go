package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("🔄 Comparing ProcessBytes vs Template Processing")
	log.Println(strings.Repeat("=", 50))

	// Read the template DOCX file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error reading template: %v", err)
	}

	// Scenario: Generate a sales report with the same data
	log.Printf("📋 Scenario: Generate sales reports for multiple customers")
	log.Println()

	customers := []map[string]interface{}{
		{
			"CompanyName": "Alpha Corp",
			"ContactName": "Alice Johnson",
			"Email":       "alice@alpha.com",
			"Orders": []map[string]interface{}{
				{"Product": "Software License", "Amount": 2500.00, "Quantity": 5},
				{"Product": "Support Package", "Amount": 800.00, "Quantity": 1},
			},
		},
		{
			"CompanyName": "Beta Industries",
			"ContactName": "Bob Chen",
			"Email":       "bob@beta.com",
			"Orders": []map[string]interface{}{
				{"Product": "Hardware System", "Amount": 5500.00, "Quantity": 1},
				{"Product": "Installation", "Amount": 1200.00, "Quantity": 1},
			},
		},
	}

	// Method 1: ProcessBytes - Manual approach (statistically impossible for multiple customers)
	log.Println("📄 Method 1: ProcessBytes Approach")
	log.Println("   ❌ Limitations:")
	log.Println("   - Requires creating separate DOCX files for each customer")
	log.Println("   - No dynamic content generation")
	log.Println("   - Manual placeholder mapping for each document")
	log.Println("   - No loops or conditional logic")
	log.Println()

	// Demonstrate the manual approach for one customer
	customer := customers[0]
	replacements := map[string]string{
		"company_name": customer["CompanyName"].(string),
		"contact_name": customer["ContactName"].(string),
		"email":        customer["Email"].(string),
		"order_total":  fmt.Sprintf("%.2f", 2500.00+800.00), // Manual calculation
	}

	result1, err := docx.ProcessBytes(templateBytes, replacements)
	if err != nil {
		log.Fatalf("❌ ProcessBytes failed: %v", err)
	}

	err = os.WriteFile("processbytes_example.docx", result1, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving ProcessBytes result: %v", err)
	}

	log.Printf("   ✅ Created single document: processbytes_example.docx")
	log.Println()

	// Method 2: Template Processing - Dynamic approach
	log.Println("🔧 Method 2: ProcessTemplate Approach")
	log.Println("   ✅ Advantages:")
	log.Println("   - Single template handles multiple customers")
	log.Println("   - Dynamic content generation with loops")
	log.Println("   - Automatic calculations and formatting")
	log.Println("   - Professional business logic")
	log.Println()

	template := `SALES REPORT - {{.Date}}

{{range .Customers}}-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
CUSTOMER: {{.CompanyName}}
CONTACT:  {{.ContactName}} ({{.Email}})

ORDER SUMMARY:
{{range .Orders}}✓ {{.Product}} ({{.Quantity|printf "%.0f"}} units) -- ${{.Amount|printf "%.2f"}}
{{end}}
Subtotal: ${{.OrdersTotal|printf "%.2f"}}
{{if .HasDiscount}}💳 Discount Applied: {{.DiscountPercent}}% ({{.DiscountAmount|printf "%.2f"}}){{end}}
Total: ${{.GrandTotal|printf "%.2f"}}

{{if .IsVIP}}⭐ VIP Customer Status{{end}}
-------------------------------------------------------------------------
{{end}}
TOTAL SALES VALUE: ${{.TotalValue|printf "%.2f"}}
TOTAL CUSTOMERS: {{len .Customers}}
AVERAGE ORDER: ${{.AverageOrder|printf "%.2f"}}

Generated on: {{.GenerationDate}}
Report ID: {{.ReportId}}`

	// Enhanced data with calculations
	reportData := map[string]interface{}{
		"Date":           "2024-01-15",
		"Customers":      customers,
		"TotalValue":     0.0,
		"AverageOrder":   0.0,
		"GenerationDate": "January 15, 2024",
		"ReportId":       "RPT-2024-001",
	}

	// Add calculations
	totalValue := 0.0
	for i, customer := range customers {
		customerOrders := customer["Orders"].([]map[string]interface{})
		customerTotal := 0.0
		for _, order := range customerOrders {
			customerTotal += order["Amount"].(float64)
		}

		// Add calculations to each customer
		customers[i]["OrdersTotal"] = customerTotal
		customers[i]["HasDiscount"] = customerTotal > 3000.0
		customers[i]["DiscountPercent"] = 5
		customers[i]["DiscountAmount"] = customerTotal * 0.05
		customers[i]["GrandTotal"] = customerTotal - (customerTotal * 0.05)
		customers[i]["IsVIP"] = customerTotal > 5000.0

		if customerTotal > 5000.0 {
			customers[i]["IsVIP"] = true
		}

		totalValue += customers[i]["GrandTotal"].(float64)
	}

	reportData["Customers"] = customers
	reportData["TotalValue"] = totalValue
	reportData["AverageOrder"] = totalValue / float64(len(customers))

	// Convert to plain text
	textReport, err := docx.ProcessTemplate(template, reportData)
	if err != nil {
		log.Fatalf("❌ ProcessTemplate failed: %v", err)
	}

	err = os.WriteFile("template_report.txt", []byte(textReport), 0644)
	if err != nil {
		log.Fatalf("❌ Error saving template result: %v", err)
	}

	log.Printf("   ✅ Generated comprehensive report: template_report.txt")
	log.Println()

	// Method 3: ProcessTemplateBytes - Best of both worlds
	log.Println("🚀 Method 3: ProcessTemplateBytes Approach")
	log.Println("   ✅ Ultimate Solution:")
	log.Println("   - Template processing + DOCX generation")
	log.Println("   - Professional document formatting")
	log.Println("   - Dynamic content in Word documents")
	log.Println("   - Maintains DOCX structure and formatting")
	log.Println()

	// Create a document-friendly template
	documentTemplate := `SALES REPORT - {{.ReportInfo.Date}}

{{range .Customers}}
===================================================================================
CUSTOMER DETAILS

Company: {{.CompanyName}}
Contact: {{.ContactName}}
Email: {{.Email}}

Orders for this quarter:
{{range .Orders}}
• {{.Product}} ({{.Quantity|printf "%.0f"}} units)
  Unit Price: ${{.Amount|printf "%.2f"}}
{{end}}

Order Total: ${{.OrdersTotal|printf "%.2f"}}
{{if .HasDiscount}}Discount Applied: {{.DiscountPercent}}% (Savings: ${{.DiscountAmount|printf "%.2f"}}){{end}}
Grand Total: ${{.GrandTotal|printf "%.2f"}}

{{if .IsVIP}}Status: ⭐ VIP Customer{{end}}
===================================================================================

{{end}}
SUMMARY STATISTICS

Total Sales Value: ${{.Summary.TotalValue|printf "%.2f"}}
Number of Customers: {{.Summary.CustomerCount}}
Average Order Value: ${{.Summary.AverageOrder|printf "%.2f"}}

{{if .Summary.HasVIPCustomers}}VIP Customers: {{.Summary.VIPCount}}{{end}}

This report was generated on {{.ReportInfo.GenerationDate}}
Report ID: {{.ReportInfo.ReportId}}
Prepared by: {{.ReportInfo.PreparedBy}}`

	// Prepare comprehensive data
	summaryData := map[string]interface{}{
		"Customers": customers,
		"ReportInfo": map[string]string{
			"Date":           "2024-01-15",
			"GenerationDate": "January 15, 2024 at 2:30 PM",
			"ReportId":       "RPT-2024-001",
			"PreparedBy":     "Sales Analytics Team",
		},
		" Summary": map[string]interface{}{
			"TotalValue":      totalValue,
			"CustomerCount":   len(customers),
			"AverageOrder":    totalValue / float64(len(customers)),
			"HasVIPCustomers": len([]bool{totalValue > 5000}) > 0,
			"VIPCount":        len(customers),
		},
	}

	// Process DOCX with template
	docConfig := docx.TemplateConfig{
		PlaceholderKey: "key",
		TemplateText:   documentTemplate,
		Data:           summaryData,
	}

	result3, err := docx.ProcessTemplateBytes(templateBytes, docConfig)
	if err != nil {
		log.Fatalf("❌ ProcessTemplateBytes failed: %v", err)
	}

	err = os.WriteFile("template_document.docx", result3, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving ProcessTemplateBytes result: %v", err)
	}

	log.Printf("   ✅ Created professional document: template_document.docx")
	log.Println()

	// Comparison Summary
	log.Println("📊 COMPARISON SUMMARY")
	log.Println(strings.Repeat("=", 25))
	log.Printf("| Method                  | Files Generated     | Use Case                          |")
	log.Println("| ---------------------- | ------------------- | --------------------------------- |")
	log.Printf("| ProcessBytes            | 1 DOCX file         | Simple placeholder replacement     |")
	log.Printf("| ProcessTemplate         | 1 TXT file         | Dynamic text generation             |")
	log.Printf("| ProcessTemplateBytes    | 1 DOCX file         | Professional document generation   |")
	log.Println()
	log.Println("🎯 RECOMMENDATIONS")
	log.Println("   • Use ProcessBytes when:")
	log.Println("     - You have simple placeholder replacement needs")
	log.Println("     - You know exactly what to replace before processing")
	log.Println("     - You're processing single documents")
	log.Println()
	log.Println("   • Use ProcessTemplate when:")
	log.Println("     - You need dynamic content generation")
	log.Println("     - You want to work with plain text directly")
	log.Println("     - You're creating reports or templated content")
	log.Println()
	log.Println("   • Use ProcessTemplateBytes when:")
	log.Println("     - You need the best of both worlds")
	log.Println("     - You're building professional document automation")
	log.Println("     - You want DOCX output with template logic")
	log.Println()
	log.Println("✨ All generated files are ready to open in Microsoft Word!")
}
