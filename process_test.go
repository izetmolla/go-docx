package docx

import (
	"os"
	"testing"
)

func TestProcessBytes(t *testing.T) {
	templateBytes, err := os.ReadFile("./test/template.docx")
	if err != nil {
		t.Fatalf("failed to read template.docx: %v", err)
	}

	replacements := map[string]string{
		"key": "REPLACED_VALUE",
	}

	outputBytes, err := ProcessBytes(templateBytes, replacements)
	if err != nil {
		t.Fatalf("ProcessBytes failed: %v", err)
	}

	// Verify the output (simplified check for now, a full docx comparison would be more robust)
	if len(outputBytes) == 0 {
		t.Errorf("expected non-empty output bytes, got empty")
	}

	// Verify output is a valid ZIP (DOCX)
	if len(outputBytes) < 4 || outputBytes[0] != 'P' || outputBytes[1] != 'K' {
		t.Errorf("ProcessBytes output is not a valid ZIP file")
	}
}

func TestProcessBytes_EmptyInput(t *testing.T) {
	replacements := map[string]string{
		"key": "value",
	}

	_, err := ProcessBytes([]byte{}, replacements)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestProcessBytes_InvalidZip(t *testing.T) {
	invalidZip := []byte("This is not a valid DOCX file")
	replacements := map[string]string{
		"key": "value",
	}

	_, err := ProcessBytes(invalidZip, replacements)
	if err == nil {
		t.Fatal("expected error for invalid ZIP input")
	}
}
