package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.adamworthington.net/internal/assert"
)

func TestPing(t *testing.T) {
    rr := httptest.NewRecorder()

    r, err := http.NewRequest(http.MethodGet, "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    ping(rr, r)

    result := rr.Result()

    assert.Equal(t, result.StatusCode, http.StatusOK)

    defer result.Body.Close()

    body, err := io.ReadAll(result.Body)
    if err != nil {
        t.Fatal(err)
    }
    body = bytes.TrimSpace(body)
    assert.Equal(t, string(body), "OK")
}
