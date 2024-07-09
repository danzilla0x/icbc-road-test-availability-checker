package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginPayload struct {
	DrvrLastName  string `json:"drvrLastName"`
	LicenceNumber string `json:"licenceNumber"`
	Keyword       string `json:"keyword"`
}

func GetBearerToken(loginUrl, lastName, licenceNumber, keyword, userAgent string) (string, error) {
	payload := LoginPayload{
		DrvrLastName:  lastName,
		LicenceNumber: licenceNumber,
		Keyword:       keyword,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", loginUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", "https://onlinebusiness.icbc.com/webdeas-ui/login;type=driver")
	req.Header.Add("Cache-Control", "no-cache, no-store")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Expires", "0")
	req.Header.Add("Sec-Ch-Ua-Platform", "macOS")
	req.Header.Add("Sec-Ch-Ua", "Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Brave\";v=\"126\"")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unsuccessful request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("bad request")
		}
		return "", fmt.Errorf("bad request: %s", body)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid credentials")
	}

	return response.Header.Get("Authorization"), nil
}
