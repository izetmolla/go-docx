package docx

import (
	"os"
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "simple substitution",
			template: "Hello {{.Name}}!",
			data:     map[string]string{"Name": "World"},
			expected: "Hello World!",
		},
		{
			name:     "multiple variables",
			template: "Dear {{.Customer}}, your order #{{.OrderID}} is {{.Status}}.",
			data: map[string]interface{}{
				"Customer": "John Doe",
				"OrderID":  "12345",
				"Status":   "shipped",
			},
			expected: "Dear John Doe, your order #12345 is shipped.",
		},
		{
			name:     "text formatting",
			template: "**{{.Title}}**\nThis is a {{.Description}}.",
			data: map[string]string{
				"Title":       "Important Notice",
				"Description": "critical update",
			},
			expected: "**Important Notice**\nThis is a critical update.",
		},
		{
			name:     "empty template",
			template: "",
			data:     map[string]string{"Name": "World"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessTemplate(tt.template, tt.data)
			if tt.template == "" {
				if err == nil {
					t.Error("Expected error for empty template")
				}
				return
			}

			if err != nil {
				t.Fatalf("ProcessTemplate failed: %v", err)
			}

			// Check that the result contains expected text content
			if tt.expected != "" && result != tt.expected {
				t.Errorf("Expected: %q, got: %q", tt.expected, result)
			}
		})
	}
}

func TestProcessTemplateBytes(t *testing.T) {
	// Load test template
	input, err := os.ReadFile("./test/template.docx")
	if err != nil {
		t.Skipf("Skipping test - no template file found: %v", err)
	}

	config := TemplateConfig{
		PlaceholderKey: "key",
		TemplateText:   "Invoice for {{.Customer}}:\nAmount: ${{.Amount}}",
		Data: map[string]interface{}{
			"Customer": "ABC Corp",
			"Amount":   "1,234.56",
		},
	}

	output, err := ProcessTemplateBytes(input, config)
	if err != nil {
		t.Fatalf("ProcessTemplateBytes failed: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("ProcessTemplateBytes returned empty output")
	}

	// Verify output is a valid ZIP (DOCX)
	if len(output) < 4 || output[0] != 'P' || output[1] != 'K' {
		t.Fatal("ProcessTemplateBytes output is not a valid ZIP file")
	}
}

func TestProcessTemplateBytes_InvalidConfig(t *testing.T) {
	input := []byte("fake docx content")

	// Test empty placeholder key
	config := TemplateConfig{
		PlaceholderKey: "",
		TemplateText:   "Hello {{.Name}}",
		Data:           map[string]string{"Name": "World"},
	}

	_, err := ProcessTemplateBytes(input, config)
	if err == nil {
		t.Fatal("Expected error for empty placeholder key")
	}

	// Test empty template
	config.PlaceholderKey = "key"
	config.TemplateText = ""

	_, err = ProcessTemplateBytes(input, config)
	if err == nil {
		t.Fatal("Expected error for empty template")
	}
}
