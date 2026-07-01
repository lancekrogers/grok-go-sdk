package grok

import (
	"context"
	"reflect"
	"testing"
)

func TestExportArgs(t *testing.T) {
	cases := []struct {
		name string
		id   string
		opts ExportOptions
		want []string
	}{
		{"base", "sess1", ExportOptions{}, []string{"export", "sess1"}},
		{"output", "sess1", ExportOptions{Output: "out.md"}, []string{"export", "sess1", "out.md"}},
		{"clipboard", "sess1", ExportOptions{Clipboard: true}, []string{"export", "sess1", "-c"}},
	}
	for _, tc := range cases {
		if got := exportArgs(tc.id, tc.opts); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %v want %v", tc.name, got, tc.want)
		}
	}
}

func TestExport_EmptyID(t *testing.T) {
	c := NewClient("/nonexistent")
	if _, err := c.Export(context.Background(), "", ExportOptions{}); err == nil {
		t.Fatal("expected error for empty sessionID")
	}
}

func TestExport_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	out, err := c.Export(context.Background(), "mock-session", ExportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected markdown output")
	}
}
