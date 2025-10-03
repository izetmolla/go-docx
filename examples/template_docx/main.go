package main

import (
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("🚀 ProcessTemplateDocx Example")
	log.Println("This demonstrates processing DOCX templates with Go template syntax")
	log.Println(strings.Repeat("=", 60))
	log.Println()

	// Read template DOCX file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error reading template file: %v", err)
	}
	log.Printf("✅ Loaded template DOCX file (%d bytes)", len(templateBytes))

	// Example 1: Simple data replacement
	log.Println()
	log.Println("📄 Example 1: Simple Template Processing")
	log.Println("   This shows how ProcessTemplateDocx handles basic {{...}} tags")

	simpleData := map[string]interface{}{
		"company": "ACME Corporation",
		"contact": "Jane Smith",
		"email":   "jane@acme.com",
		"phone":   "(555) 123-4567",
		"address": "123 Business Ave\nSuite 200\nNew York, NY 10001",
		"website": "www.acme.com",
	}

	// Process DOCX with simple template data
	result1, err := docx.ProcessTemplateDocx(templateBytes, simpleData)
	if err != nil {
		log.Fatalf("❌ Simple template processing failed: %v", err)
	}

	// Save result
	err = os.WriteFile("simple_template_processed.docx", result1, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving simple processed file: %v", err)
	}
	log.Printf("✅ Created: simple_template_processed.docx")

	// Example 2: Advanced template with conditions and loops
	log.Println()
	log.Println("🔧 Example 2: Advanced Template Processing")
	log.Println("   This demonstrates complex template logic with {{if}} and {{range}}")

	advancedData := map[string]interface{}{
		"invoiceNumber": "INV-2024-001",
		"customer": map[string]interface{}{
			"name":    "Tech Solutions LLC",
			"address": "456 Innovation Drive\nAustin, TX 78701",
			"email":   "billing@techsolutions.com",
			"phone":   "(555) 987-6543",
			"isVip":   true,
		},
		"date": "January 15, 2024",
		"items": []map[string]interface{}{
			{"product": "Enterprise Software License", "quantity": 5, "price": 2500.00},
			{"product": "Professional Services (100hrs)", "quantity": 1, "price": 150.00},
			{"product": "Premium Support Package", "quantity": 1, "price": 2000.00},
		},
		"subtotal":     29500.00,
		"discount":     5,
		"tax":          2101.25,
		"shipping":     100.00,
		"total":        30226.25,
		"paymentTerms": "Net 30",
		"dueDate":      "February 15, 2024",
	}

	// Process DOCX with advanced template data
	result2, err := docx.ProcessTemplateDocx(templateBytes, advancedData)
	if err != nil {
		log.Fatalf("❌ Advanced template processing failed: %v", err)
	}

	// Save result
	err = os.WriteFile("advanced_template_processed.docx", result2, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving advanced processed file: %v", err)
	}
	log.Printf("✅ Created: advanced_template_processed.docx")

	// Example 3: Business report
	log.Println()
	log.Println("📊 Example 3: Business Report Processing")
	log.Println("   This shows template processing for complex business documents")

	reportData := map[string]interface{}{
		"reportDate":    "Q1 2024",
		"companyName":   "Global Analytics Inc",
		"preparedBy":    "Accounting Department",
		"quarterEnding": "March 31, 2024",
		"revenue": map[string]interface{}{
			"product":  150000.00,
			"services": 75000.00,
			"support":  25000.00,
			"total":    250000.00,
		},
		"expenses": map[string]interface{}{
			"salaries":   120000.00,
			"marketing":  15000.00,
			"operations": 8000.00,
			"total":      143000.00,
		},
		"netIncome":   107000.00,
		"grossMargin": 62.4, // percentage
		"topClients": []map[string]interface{}{
			{"name": "Alpha Corp", "amount": 45000.00, "status": "Paid"},
			{"name": "Beta Industries", "amount": 32000.00, "status": "Pending"},
			{"name": "Gamma Solutions", "amount": 28000.00, "status": "Paid"},
		},
		"notes": "Q1 performance exceeded expectations with strong growth in services revenue.",
	}

	// Process DOCX with report data
	result3, err := docx.ProcessTemplateDocx(templateBytes, reportData)
	if err != nil {
		log.Fatalf("❌ Report processing failed: %v", err)
	}

	// Save result
	err = os.WriteFile("business_report.docx", result3, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving business report: %v", err)
	}
	log.Printf("✅ Created: business_report.docx")

	// Summary
	log.Println()
	log.Println(strings.Repeat("=", 60))
	log.Println("📋 SUMMARY")
	log.Println(strings.Repeat("=", 60))
	log.Printf("✅ ProcessTemplateDocx processed 3 different scenarios")
	log.Printf("📁 Generated files:")
	log.Printf("   • simple_template_processed.docx    - Basic template processing")
	log.Printf("   • advanced_template_processed.docx - Complex business logic")
	log.Printf("   • business_report.docx            - Financial report")
	log.Println()
	log.Println("🎯 KEY FEATURES:")
	log.Println("   • Direct DOCX template processing")
	log.Println("   • Full Go template syntax support")
	log.Println("   • Conditions, loops, and formatting")
	log.Println("   • Maintains Word document structure")
	log.Println("   • No HTML/RTF conversion needed")
	log.Println()
	log.Println("💡 ProcessTemplateDocx vs ProcessTemplate:")
	log.Println("   • ProcessTemplateDocx: DOCX template → DOCX output")
	log.Println("   • ProcessTemplate: String template → Plain text output")
	log.Println("   • Both support full Go template syntax")
	log.Println()
	log.Println("✨ All files are ready to open in Microsoft Word!")
}
