package docx

import (
	"os"
	"testing"
)

func TestProcessTemplateDocx(t *testing.T) {
	tests := []struct {
		name        string
		templateDoc string // Path to template DOCX file
		data        interface{}
		expectError bool
	}{
		{
			name:        "simple template replacement",
			data:        map[string]string{"customer": "ABC Corp", "amount": "100.00"},
			expectError: false,
		},
		{
			name:        "complex template with struct",
			data:        map[string]interface{}{"company": "Test Inc", "isVip": true},
			expectError: false,
		},
		{
			name:        "empty input",
			data:        map[string]string{},
			expectError: true,
		},
	}

	// Read test template file if available
	input, err := os.ReadFile("./test/template.docx")
	if err != nil {
		t.Skipf("Skipping test - no template file found: %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testInput []byte
			if tt.expectError && tt.name == "empty input" {
				testInput = []byte{} // Empty input for this test case
			} else {
				testInput = input // Use valid DOCX for other tests
			}

			result, err := ProcessTemplateDocx(testInput, tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("ProcessTemplateDocx failed: %v", err)
			}

			if len(result) == 0 {
				t.Error("Expected non-empty result")
			}

			// Verify result is a valid ZIP (DOCX)
			if len(result) < 4 || result[0] != 'P' || result[1] != 'K' {
				t.Error("Result is not a valid ZIP file")
			}
		})
	}
}

func TestProcessTemplateDocx_TemplateSyntax(t *testing.T) {
	// This test would require a template DOCX file with {{...}} syntax
	// For now, we'll test the helper functions

	t.Run("isTemplatePlaceholder", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"{{.Name}}", true},
			{"{{if .Condition}}", true},
			{"{{range .Items}}", true},
			{"Hello World", false},
			{"{placeholder}", false},
			{"{{.Name}", false}, // Incomplete
			{".Name}}", false},  // Incomplete
		}

		for _, tt := range tests {
			result := isTemplatePlaceholder(tt.input)
			if result != tt.expected {
				t.Errorf("isTemplatePlaceholder(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("executeTemplatePlaceholder", func(t *testing.T) {
		template := "Hello {{.Name}}!"
		data := map[string]string{"Name": "World"}

		result, err := executeTemplatePlaceholder(template, data)
		if err != nil {
			t.Fatalf("executeTemplatePlaceholder failed: %v", err)
		}

		expected := "Hello World!"
		if result != expected {
			t.Errorf("executeTemplatePlaceholder result = %s, expected %s", result, expected)
		}
	})

	t.Run("executeTemplatePlaceholder complex", func(t *testing.T) {
		template := "Order #{{.OrderID}}: ${{printf \"%.2f\" .Amount}}"
		data := map[string]interface{}{
			"OrderID": "12345",
			"Amount":  99.99,
		}

		result, err := executeTemplatePlaceholder(template, data)
		if err != nil {
			t.Fatalf("executeTemplatePlaceholder failed: %v", err)
		}

		expected := "Order #12345: $99.99"
		if result != expected {
			t.Errorf("executeTemplatePlaceholder result = %s, expected %s", result, expected)
		}
	})
}

func TestProcessTemplateDocx_InvalidTemplate(t *testing.T) {
	input := []byte("fake docx")
	data := map[string]string{"key": "value"}

	_, err := ProcessTemplateDocx(input, data)
	if err == nil {
		t.Error("Expected error for invalid DOCX input")
	}
}

func TestProcessTemplateDocx_EmptyData(t *testing.T) {
	input, err := os.ReadFile("./test/template.docx")
	if err != nil {
		t.Skipf("Skipping test - no template file found: %v", err)
		return
	}

	_, err = ProcessTemplateDocx(input, nil)
	if err != nil {
		t.Fatalf("ProcessTemplateDocx should handle nil data gracefully")
	}
}

func TestExecuteTemplatePlaceholder_MissingData(t *testing.T) {
	t.Run("missing data preserves template", func(t *testing.T) {
		template := "Hello {{.MissingField}}, order #{{.ExistingField}}"
		data := map[string]interface{}{
			"ExistingField": "12345",
			// MissingField is intentionally missing
		}

		result, err := executeTemplatePlaceholder(template, data)
		if err != nil {
			t.Fatalf("executeTemplatePlaceholder failed: %v", err)
		}

		// Should preserve original template since MissingField is missing
		if result != template {
			t.Errorf("Expected original template to be preserved, got: %s", result)
		}
	})

	t.Run("complete data replaces template", func(t *testing.T) {
		template := "Hello {{.Name}}, order #{{.OrderID}}"
		data := map[string]interface{}{
			"Name":    "John",
			"OrderID": "12345",
		}

		result, err := executeTemplatePlaceholder(template, data)
		if err != nil {
			t.Fatalf("executeTemplatePlaceholder failed: %v", err)
		}

		expected := "Hello John, order #12345"
		if result != expected {
			t.Errorf("Expected %s, got: %s", expected, result)
		}
	})

	t.Run("partial data with missing field", func(t *testing.T) {
		template := "Hello {{.Name}}, {{.MissingField}} - Order: {{.OrderID}}"
		data := map[string]interface{}{
			"Name":    "John",
			"OrderID": "12345",
			// MissingField is missing
		}

		result, err := executeTemplatePlaceholder(template, data)
		if err != nil {
			t.Fatalf("executeTemplatePlaceholder failed: %v", err)
		}

		// Should preserve original template since MissingField is missing
		if result != template {
			t.Errorf("Expected original template to be preserved, got: %s", result)
		}
	})
}
