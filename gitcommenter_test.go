package gitcommenter

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.OllamaEndpoint != "http://localhost:11434" {
		t.Errorf("Expected default endpoint to be http://localhost:11434, got %s", config.OllamaEndpoint)
	}

	if config.Model != "llama2" {
		t.Errorf("Expected default model to be llama2, got %s", config.Model)
	}

	if config.Temperature != 0.7 {
		t.Errorf("Expected default temperature to be 0.7, got %f", config.Temperature)
	}
}

func TestNew(t *testing.T) {
	config := DefaultConfig()
	commenter := New(config)

	if commenter == nil {
		t.Error("Expected New to return a non-nil GitCommenter")
	}

	if commenter.config != config {
		t.Error("Expected GitCommenter to use the provided config")
	}
}

func TestParseChangeType(t *testing.T) {
	commenter := New(nil)

	tests := []struct {
		status   string
		expected string
	}{
		{"A", "added"},
		{"M", "modified"},
		{"D", "deleted"},
		{"R", "renamed"},
		{"C", "copied"},
		{"?", "modified"}, // default case
	}

	for _, test := range tests {
		result := commenter.parseChangeType(test.status)
		if result != test.expected {
			t.Errorf("parseChangeType(%s) = %s, want %s", test.status, result, test.expected)
		}
	}
}

func TestCountDiffLines(t *testing.T) {
	commenter := New(nil)

	diff := `--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,4 @@
 line1
+added line
 line2
-removed line
 line3
+another added line`

	added, removed := commenter.countDiffLines(diff)

	if added != 2 {
		t.Errorf("Expected 2 added lines, got %d", added)
	}

	if removed != 1 {
		t.Errorf("Expected 1 removed line, got %d", removed)
	}
}

func TestBuildChangeContext(t *testing.T) {
	commenter := New(nil)

	changes := []FileChange{
		{
			FilePath:     "file1.go",
			ChangeType:   "modified",
			LinesAdded:   5,
			LinesRemoved: 2,
		},
		{
			FilePath:     "file2.txt",
			ChangeType:   "added",
			LinesAdded:   10,
			LinesRemoved: 0,
		},
	}

	context := commenter.buildChangeContext(changes)

	if !contains(context, "file1.go") {
		t.Error("Expected context to contain file1.go")
	}

	if !contains(context, "file2.txt") {
		t.Error("Expected context to contain file2.txt")
	}

	if !contains(context, "modified") {
		t.Error("Expected context to contain 'modified'")
	}

	if !contains(context, "added") {
		t.Error("Expected context to contain 'added'")
	}
}

func TestParseCommitSuggestion(t *testing.T) {
	commenter := New(nil)

	changes := []FileChange{
		{FilePath: "file1.go"},
		{FilePath: "file2.txt"},
	}

	response := "feat: add new functionality\n\nThis commit adds new features to improve user experience."

	suggestion := commenter.parseCommitSuggestion(response, changes)

	if suggestion.Subject != "feat: add new functionality" {
		t.Errorf("Expected subject 'feat: add new functionality', got '%s'", suggestion.Subject)
	}

	expectedBody := "This commit adds new features to improve user experience."
	if suggestion.Body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, suggestion.Body)
	}

	if len(suggestion.FilesAffected) != 2 {
		t.Errorf("Expected 2 affected files, got %d", len(suggestion.FilesAffected))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
