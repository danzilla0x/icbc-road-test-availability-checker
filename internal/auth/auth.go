package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	LOGIN_URL = "https://onlinebusiness.icbc.com/deas-api/v1/webLogin/webLogin"
)

type LoginPayload struct {
	DrvrLastName  string `json:"drvrLastName"`
	LicenceNumber string `json:"licenceNumber"`
	Keyword       string `json:"keyword"`
}

func GetBearerToken(lastName, licenceNumber, keyword, userAgent string) (string, error) {
	payload := LoginPayload{
		DrvrLastName:  lastName,
		LicenceNumber: licenceNumber,
		Keyword:       keyword,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", LOGIN_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", "https://onlinebusiness.icbc.com/webdeas-ui/login;type=driver")
	req.Header.Add("Cache-Control", "no-cache, no-store")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unsuccessful request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("invalid credentials")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s, status code: %d", response.Status, response.StatusCode)
	}

	return response.Header.Get("Authorization"), nil
}
