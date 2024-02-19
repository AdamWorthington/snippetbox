package assert

import (
    "testing"
    "strings"
)

func NilError(t *testing.T, actual error) {
    t.Helper()

    if actual != nil {
        t.Errorf("got %v, expected nil", actual)
    }
}

func StringContains(t *testing.T, str, substr string) {
    t.Helper()

    if !strings.Contains(str, substr) {
        t.Errorf("%v not found in: %v", substr, str)
    }
}

func Equal[T comparable](t *testing.T, actual, expected T) {
    t.Helper()

    if actual != expected {
        t.Errorf("got %v; want %v", actual, expected)
    }
}
