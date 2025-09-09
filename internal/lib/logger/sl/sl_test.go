package sl

import (
	"errors"
	"log/slog"
	"testing"
)

func TestErr(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		attr := Err(nil)

		if attr.Key != "error" {
			t.Errorf("Expected key 'error', got %q", attr.Key)
		}

		if attr.Value.Kind() != slog.KindString {
			t.Errorf("Expected String value, got %v", attr.Value.Kind())
		}

		if attr.Value.String() != "nil" {
			t.Errorf("Expected value 'nil', got %q", attr.Value.String())
		}
	})

	t.Run("with actual error", func(t *testing.T) {
		testError := errors.New("database connection failed")
		attr := Err(testError)

		if attr.Key != "error" {
			t.Errorf("Expected key 'error', got %q", attr.Key)
		}

		if attr.Value.Kind() != slog.KindString {
			t.Errorf("Expected String value, got %v", attr.Value.Kind())
		}

		if attr.Value.String() != "database connection failed" {
			t.Errorf("Expected value 'database connection failed', got %q", attr.Value.String())
		}
	})

	t.Run("with empty error", func(t *testing.T) {
		testError := errors.New("")
		attr := Err(testError)

		if attr.Value.String() != "" {
			t.Errorf("Expected empty string, got %q", attr.Value.String())
		}
	})
}
