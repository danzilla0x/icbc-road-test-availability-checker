package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	t.Run("Invalid creds", func(t *testing.T) {
		// Create a mock server
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "One or more of your entries is incorrect. Verify that your last name, BC driver’s licence number, and keyword are correct"}`))
		}))
		defer mockServer.Close()

		_, err := GetBearerToken(mockServer.URL, "Smith", "12356", "secret", "")

		if err == nil {
			t.Errorf("expected error: %v, got: %v", true, err)
		}
	})
}
