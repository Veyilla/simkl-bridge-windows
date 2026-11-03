package desktop

import (
	"os"
	"testing"
)

func TestDocumentReadersExistOnWindows(t *testing.T) {
	for _, name := range []string{"LICENSE", "THIRD_PARTY_NOTICES.md", "licenses"} {
		if _, err := os.Stat(documentProgram(name)); err != nil {
			t.Fatalf("Windows-Leseprogramm für %s fehlt: %v", name, err)
		}
	}
}
