package ui

import (
	"strings"
	"testing"
)

func TestDialogView(t *testing.T) {
	title := "Confirmation"
	message := "Are you sure you want to continue?"
	
	dialog := NewDialog(title, message)
	result := dialog.View()
	
	// Check if the result contains the title and message
	if !strings.Contains(result, title) {
		t.Errorf("Expected view to contain title %q, but it didn't.\nGot: %s", title, result)
	}
	
	if !strings.Contains(result, message) {
		t.Errorf("Expected view to contain message %q, but it didn't.\nGot: %s", message, result)
	}
	
	// Check if the result contains the buttons
	if !strings.Contains(result, "Yes") {
		t.Errorf("Expected view to contain 'Yes' button, but it didn't.\nGot: %s", result)
	}
	
	if !strings.Contains(result, "No") {
		t.Errorf("Expected view to contain 'No' button, but it didn't.\nGot: %s", result)
	}
}
