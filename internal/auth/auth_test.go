package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBearerToken(t *testing.T) {

	testCases := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedError  bool
	}{
		{
			name:           "Successful Login",
			responseStatus: http.StatusOK,
			responseBody:   `{"Authorization": "Bearer some-token"}`,
			expectedError:  false,
		},
		{
			name:           "Status Bad Request",
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"error": "One or more of your entries is incorrect. Verify that your last name, BC driver’s licence number, and keyword are correct"}`,
			expectedError:  true,
		},
		{
			name:           "Server Error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "server error"}`,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.responseBody))
			}))
			defer mockServer.Close()

			_, err := GetBearerToken(mockServer.URL, "Smith", "12356", "secret", "")

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
