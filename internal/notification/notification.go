package notification

import (
	"net/http"
	"net/url"
)

func SendNotification(userKey, token, message string) error {
	_, err := http.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {token},
		"user":    {userKey},
		"message": {"There is an available date: " + message},
	})

	return err
}
