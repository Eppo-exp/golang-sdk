package eppoclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`OK`))
		case "/unauthorized":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`Unauthorized`))
		case "/internal-error":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`Internal Server Error`))
		case "/bad-response":
			w.WriteHeader(http.StatusOK)
			if hijacker, ok := w.(http.Hijacker); ok {
				conn, _, _ := hijacker.Hijack()
				conn.Close() // Close the connection to simulate an unreadable body
			}
		}
	}))
	defer server.Close()

	client := &http.Client{}
	hc := newHttpClient(server.URL, client, SDKParams{
		sdkKey:     "testSdkKey",
		sdkName:    "testSdkName",
		sdkVersion: "testSdkVersion",
	})

	tests := []struct {
		name           string
		resource       string
		expectedError  string
		expectedResult string
	}{
		{
			name:           "api returns http 200",
			resource:       "/test",
			expectedResult: "OK",
		},
		{
			name:          "api returns 401 unauthorized error",
			resource:      "/unauthorized",
			expectedError: "unauthorized access",
		},
		{
			name:          "api returns an 500 error",
			resource:      "/internal-error",
			expectedError: "server error: 500",
		},
		{
			name:          "api returns unreadable body",
			resource:      "/bad-response",
			expectedError: "server error: unreadable body",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := hc.get(tc.resource)
			if err != nil {
				if err.Error() != tc.expectedError {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
				if result != "" { // Check if result is not an empty string when an error is expected
					t.Errorf("Expected result to be an empty string when there is an error, got %v", result)
				}
			} else {
				if tc.expectedError != "" {
					t.Errorf("Expected error %v, got nil", tc.expectedError)
				}
				if result != tc.expectedResult {
					t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
				}
			}
		})
	}
}
