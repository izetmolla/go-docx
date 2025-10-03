package docx

import (
	"strings"
	"testing"
)

func TestAssembleFullPlaceholdersWithMismatchedArrays(t *testing.T) {
	run := &Run{ID: 1}

	// Test case 1: more openPos than closePos
	openPos1 := []int{0, 10, 20}
	closePos1 := []int{8, 18} // Only 2 close positions

	placeholders1 := assembleFullPlaceholders(run, openPos1, closePos1)
	if len(placeholders1) != 2 {
		t.Errorf("Expected 2 placeholders for mismatched arrays (3 opens, 2 closes), got %d", len(placeholders1))
	}

	// Test case 2: more closePos than openPos
	openPos2 := []int{0, 10}
	closePos2 := []int{8, 18, 30} // 3 close positions, only 2 opens

	placeholders2 := assembleFullPlaceholders(run, openPos2, closePos2)
	if len(placeholders2) != 2 {
		t.Errorf("Expected 2 placeholders for mismatched arrays (2 opens, 3 closes), got %d", len(placeholders2))
	}

	// Test case 3: empty arrays
	placeholders3 := assembleFullPlaceholders(run, []int{}, []int{})
	if len(placeholders3) != 0 {
		t.Errorf("Expected 0 placeholders for empty arrays, got %d", len(placeholders3))
	}

	// Test case 4: equal length arrays (normal case)
	openPos4 := []int{0, 10}
	closePos4 := []int{8, 18}

	placeholders4 := assembleFullPlaceholders(run, openPos4, closePos4)
	if len(placeholders4) != 2 {
		t.Errorf("Expected 2 placeholders for equal length arrays, got %d", len(placeholders4))
	}
}

func TestNestedTemplateProcessing(t *testing.T) {
	// Test the ProcessTemplate function with nested conditionals that could cause issues
	templateText := `{{- if eq .Name "Izet" -}}
Happy Birthday!
{{- else -}}
Hello {{.Name}}.
{{- end -}}
Best regards, {{.unknown_field}}`

	data := map[string]interface{}{
		"Name": "Izet",
		// Intentionally missing "unknown_field"
	}

	result, err := ProcessTemplate(templateText, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	// Verify that the unknown field is handled (outputs <no value> as per Go's text/template behavior)
	if !strings.Contains(result, "<no value>") {
		t.Errorf("Expected template to handle unknown field with <no value>, got: %s", result)
	}

	// Verify the conditional logic worked
	if !strings.Contains(result, "Happy Birthday!") {
		t.Errorf("Expected conditional to evaluate correctly, got: %s", result)
	}
}
