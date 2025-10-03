package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("🔄 Comparing DOCX Processing Methods")
	log.Println(strings.Repeat("=", 40))

	// Read the template DOCX file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error reading template: %v", err)
	}

	// Scenario: Generate reports for multiple customers using different approaches
	log.Printf("📋 Scenario: Generate sales reports for multiple customers")
	log.Println()

	customers := []map[string]interface{}{
		{
			"company_name": "Alpha Corp",
			"contact_name": "Alice Johnson",
			"email":        "alice@alpha.com",
			"order_total":  "3,300.00",
		},
		{
			"company_name": "Beta Industries",
			"contact_name": "Bob Chen", 
			"email":        "bob@beta.com",
			"order_total":  "6,700.00",
		},
	}

	// Method 1: Document API - Full Control
	log.Println("📄 Method 1: Document API - Full Control")
	log.Println("   ✅ Advantages:")
	log.Println("   - Complete control over document processing")
	log.Println("   - Can access individual placeholders")
	log.Println("   - Can modify multiple files (headers, footers)")
	log.Println("   - Memory efficient for large documents")
	log.Println()

	for i, customer := range customers {
		doc, err := docx.OpenBytes(templateBytes)
		if err != nil {
			log.Fatalf("❌ Document API failed: %v", err)
		}

		// Replace placeholders with customer data
		doc.ReplaceAll(docx.PlaceholderMap{
			"company":     customer["company_name"].(string),
			"contact":     customer["contact_name"].(string),
			"email":       customer["email"].(string),
			"order_total": customer["order_total"].(string),
		})

		filename := fmt.Sprintf("document_api_customer_%d.docx", i+1)
		err = doc.WriteToFile(filename)
		doc.Close()
		if err != nil {
			log.Fatalf("❌ Error saving document: %v", err)
		}
		log.Printf("   ✅ Created: %s", filename)
	}

	log.Println()

	// Method 2: ProcessBytes - Batch Processing
	log.Println("🚀 Method 2: ProcessBytes - Batch Processing")
	log.Println("   ✅ Advantages:")
	log.Println("   - Fast in-memory processing")
	log.Println("   - Simple API for bulk operations")
	log.Println("   - Great for server-side processing")
	log.Println("   - No file handle management needed")
	log.Println()

	for i, customer := range customers {
		replacements := map[string]string{
			"company":     customer["company_name"].(string),
			"contact":     customer["contact_name"].(string),
			"email":       customer["email"].(string),
			"order_total": customer["order_total"].(string),
		}

		outputBytes, err := docx.ProcessBytes(templateBytes, replacements)
		if err != nil {
			log.Fatalf("❌ ProcessBytes failed: %v", err)
		}

		filename := fmt.Sprintf("processbytes_customer_%d.docx", i+1)
		err = os.WriteFile(filename, outputBytes, 0644)
		if err != nil {
			log.Fatalf("❌ Error saving file: %v", err)
		}
		log.Printf("   ✅ Created: %s", filename)
	}

	log.Println()

	// Method 3: Hybrid approach - Multiple operations
	log.Println("🔧 Method 3: Document API with Multiple Operations")
	log.Println("   ✅ Advanced Features:")
	log.Println("   - Replace individual placeholders")
	log.Println("   - Chain multiple operations")
	log.Println("   - Access document metadata")
	log.Println("   - Custom processing logic")
	log.Println()

	doc, err := docx.OpenBytes(templateBytes)
	if err != nil {
		log.Fatalf("❌ Document API failed: %v", err)
	}
	defer doc.Close()

	// Process first customer with individual replacements
	customer1 := customers[0]
	doc.Replace("company", customer1["company_name"].(string))
	doc.Replace("contact", customer1["contact_name"].(string))
	doc.Replace("email", customer1["email"].(string))
	doc.Replace("order_total", customer1["order_total"].(string))

	err = doc.WriteToFile("hybrid_approach.docx")
	if err != nil {
		log.Fatalf("❌ Error saving hybrid document: %v", err)
	}
	log.Printf("   ✅ Created: hybrid_approach.docx")

	log.Println()

	// Comparison Summary
	log.Println("📊 COMPARISON SUMMARY")
	log.Println(strings.Repeat("=", 25))
	log.Printf("| Method                  | Files Generated     | Use Case                          |")
	log.Println("| ---------------------- | ------------------- | --------------------------------- |")
	log.Printf("| Document API            | Multiple DOCX files | Full control, complex operations   |")
	log.Printf("| ProcessBytes            | Multiple DOCX files | Fast batch processing             |")
	log.Printf("| Hybrid Approach         | Single DOCX file    | Custom processing logic           |")
	log.Println()

	log.Println("🎯 RECOMMENDATIONS")
	log.Println("   • Use Document API when:")
	log.Println("     - You need fine-grained control")
	log.Println("     - Working with multiple document files")
	log.Println("     - Complex document operations")
	log.Println()
	log.Println("   • Use ProcessBytes when:")
	log.Println("     - You need fast batch processing")
	log.Println("     - Simple placeholder replacement")
	log.Println("     - Server-side document generation")
	log.Println()
	log.Println("   • Use Hybrid approaches when:")
	log.Println("     - You need custom processing logic")
	log.Println("     - Chaining multiple operations")
	log.Println("     - Dynamically choosing operations")
	log.Println()
	log.Println("✨ All generated files are ready to open in Microsoft Word!")
}