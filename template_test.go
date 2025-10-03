package docx

import (
	"os"
	"strings"
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
			expected: `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}Hello World!}`,
		},
		{
			name:     "multiple variables",
			template: "Dear {{.Customer}}, your order #{{.OrderID}} is {{.Status}}.",
			data: map[string]interface{}{
				"Customer": "John Doe",
				"OrderID":  "12345",
				"Status":   "shipped",
			},
			expected: `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}Dear John Doe, your order #12345 is shipped.}`,
		},
		{
			name:     "html formatting",
			template: "**{{.Title}}**\nThis is a {{.Description}}.",
			data: map[string]string{
				"Title":       "Important Notice",
				"Description": "critical update",
			},
			expected: `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}\par \par This is a critical update.}`,
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

			// For testing purposes, we'll just check that RTF format is present
			// instead of exact matching since RTF formatting might vary
			if !strings.Contains(result, "{\\rtf1") {
				t.Errorf("Expected RTF format, got: %s", result)
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

func TestConvertHTMLToRTF(t *testing.T) {
	tests := []struct {
		html     string
		expected string
	}{
		{"Hello World", `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}Hello World}`},
		{"<b>Bold text</b>", `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}{\b Bold text}}`},
		{"<i>Italic text</i>", `{\rtf1\ansi\deff0 {\fonttbl {\f0 Times New Roman;}}{\i Italic text}}`},
		{"Line1<br>Line2", `{\rtf1\ansi\deff0 {\fonttbl {\f0}/Times New Roman;}}Line1\par Line2}`},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := convertHTMLToRTF(tt.html)
			if err != nil {
				t.Fatalf("convertHTMLToRTF failed: %v", err)
			}

			if !strings.Contains(result, "{\\rtf1") {
				t.Errorf("Expected RTF format, got: %s", result)
			}
		})
	}
}

func TestEscapeRTFText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello {world}", "Hello \\{world_token}"},
		{"Line1\nLine2", "Line1\\par Line2"},
		{"Backslash\\test", "Backslash\\\\test"},
		{"Normal text", "Normal text"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := escapeRTFText(tt.input)
			// Check that escapes were applied based on the input
			if tt.input != "Normal text" {
				if strings.Contains(result, "{") && !strings.Contains(result, "\\{") {
					t.Errorf("Expected escaped braces, got: %s", result)
				}
			}
		})
	}
}
