package grok

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSessionsList_AgainstFixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "sessions-list.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("fixture not yet captured")
	}
	got := parseSessions(data)
	if len(got) == 0 {
		t.Fatal("no sessions parsed")
	}
	if got[0].ID == "" {
		t.Fatal("first session missing id")
	}
}
