package main

import (
	"log"
	"os"
	"strings"

	"github.com/izetmolla/go-docx"
)

func main() {
	log.Println("🧪 ProcessTemplateDocx - Missing Data Handling")
	log.Println("This demonstrates how missing data preserves original {{...}} tags")
	log.Println(strings.Repeat("=", 65))
	log.Println()

	// Read template DOCX file
	templateBytes, err := os.ReadFile("../../test/template.docx")
	if err != nil {
		log.Fatalf("❌ Error reading template file: %v", err)
	}
	log.Printf("✅ Loaded template DOCX file (%d bytes)", len(templateBytes))

	// Example 1: Complete data - all fields present
	log.Println()
	log.Println("📄 Example 1: Complete Data (All Fields Present)")
	log.Println("   This shows normal processing when all required data is available")

	completeData := map[string]interface{}{
		"customerName": "Alice Johnson",
		"company":      "Tech Corp Inc",
		"email":        "alice@techcorp.com",
		"orderNumber":  "ORD-2024-001",
		"amount":       1250.00,
		"status":       "completed",
		"dueDate":      "2024-02-15",
		"notes":        "All information provided",
	}

	result1, err := docx.ProcessTemplateDocx(templateBytes, completeData)
	if err != nil {
		log.Fatalf("❌ Complete data processing failed: %v", err)
	}

	err = os.WriteFile("complete_data_output.docx", result1, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving complete data result: %v", err)
	}
	log.Printf("✅ Created: complete_data_output.docx (all fields replaced)")

	// Example 2: Partial data - some fields missing
	log.Println()
	log.Println("⚠️  Example 2: Partial Data (Some Fields Missing)")
	log.Println("   This shows how missing fields preserve original {{...}} tags")

	partialData := map[string]interface{}{
		"customerName": "Bob Smith",
		"company":      "Beta Industries",
		// email is missing
		"orderNumber": "ORD-2024-002",
		// amount is missing
		// status is missing
		// dueDate is missing
		"notes": "Some information missing",
	}

	result2, err := docx.ProcessTemplateDocx(templateBytes, partialData)
	if err != nil {
		log.Fatalf("❌ Partial data processing failed: %v", err)
	}

	err = os.WriteFile("partial_data_output.docx", result2, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving partial data result: %v", err)
	}
	log.Printf("✅ Created: partial_data_output.docx (missing fields preserved as {{...}})")

	// Example 3: Minimal data - most fields missing
	log.Println()
	log.Println("🔍 Example 3: Minimal Data (Most Fields Missing)")
	log.Println("   This shows error handling when very little data is provided")

	minimalData := map[string]interface{}{
		"customerName": "Carol Davis",
		// Most other fields are missing
	}

	result3, err := docx.ProcessTemplateDocx(templateBytes, minimalData)
	if err != nil {
		log.Fatalf("❌ Minimal data processing failed: %v", err)
	}

	err = os.WriteFile("minimal_data_output.docx", result3, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving minimal data result: %v", err)
	}
	log.Printf("✅ Created: minimal_data_output.docx (most fields preserved as {{...}})")

	// Example 4: No data at all
	log.Println()
	log.Println("🚫 Example 4: No Data Provided")
	log.Println("   This shows behavior with completely empty data")

	emptyData := map[string]interface{}{}

	result4, err := docx.ProcessTemplateDocx(templateBytes, emptyData)
	if err != nil {
		log.Fatalf("❌ Empty data processing failed: %v", err)
	}

	err = os.WriteFile("no_data_output.docx", result4, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving empty data result: %v", err)
	}
	log.Printf("✅ Created: no_data_output.docx (all fields preserved as {{...}})")

	// Example 5: Complex template with conditionals
	log.Println()
	log.Println("🎯 Example 5: Complex Template with Conditionals")
	log.Println("   This shows missing data handling in conditional templates")

	conditionalData := map[string]interface{}{
		"customerName": "David Wilson",
		"company":      "Delta Corp",
		"isVipCustomer": true, // Present
		// vipDiscountPercentage is missing
		// additionalServices is missing
		"standardServices": []string{"Support", "Updates"},
	}

	result5, err := docx.ProcessTemplateDocx(templateBytes, conditionalData)
	if err != nil {
		log.Fatalf("❌ Conditional data processing failed: %v", err)
	}

	err = os.WriteFile("conditional_data_output.docx", result5, 0644)
	if err != nil {
		log.Fatalf("❌ Error saving conditional data result: %v", err)
	}
	log.Printf("✅ Created: conditional_data_output.docx (some conditionals preserved)")

	// Summary
	log.Println()
	log.Println(strings.Repeat("=", 65))
	log.Println("📋 MISSING DATA HANDLING SUMMARY")
	log.Println(strings.Repeat("=", 65))
	log.Printf("📁 Generated files:")
	log.Printf("   • complete_data_output.docx      - All fields replaced")
	log.Printf("   • partial_data_output.docx      - Some {{...}} tags preserved")
	log.Printf("   • minimal_data_output.docx      - Most {{...}} tags preserved")
	log.Printf("   • no_data_output.docx           - All {{...}} tags preserved")
	log.Printf("   • conditional_data_output.docx  - Conditional {{...}} preserved")
	log.Println()
	log.Println("🎯 KEY BEHAVIOR:")
	log.Println("   ✅ Present data → Values are replaced")
	log.Println("   ✅ Missing data → Original {{...}} tags are preserved")
	log.Println("   ✅ Partial data → Only available fields are replaced")
	log.Println("   ✅ Empty data → All {{...}} tags remain unchanged")
	log.Println()
	log.Println("💡 Benefits:")
	log.Println("   • Easy to identify missing fields in output documents")
	log.Println("   • Maintains template structure for future processing")
	log.Println("   • No dangerous null/empty replacements")
	log.Println("   • Clear visual indication of incomplete data")
	log.Println()
	log.Println("🔍 Usage Scenario:")
	log.Println("   1. Create DOCX template with {{.FieldName}} tags")
	log.Println("   2. Process with incomplete data")
	log.Println("   3. Review output to see which {{...}} tags remain")
	log.Println("   4. Collect missing data and reprocess")
	log.Println()
	log.Println("✨ Open the generated documents to see {{...}} preservation in action!")
}
