package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkMultiThreadedHandler(b *testing.B) {
	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Create a multi-threaded handler with desired settings
	handler := NewMultiThreadedHandler(5, 10)

	// Reset the benchmark timer
	b.ResetTimer()

	// Perform the benchmarking
	for i := 0; i < b.N; i++ {
		// Serve the test request using the multi-threaded handler
		handler.ServeHTTP(w, req)
	}
}
