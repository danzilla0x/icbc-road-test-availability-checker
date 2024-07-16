package notification

import (
	"net/http"
	"net/url"
)

const (
	pushoverUrl  = "https://api.pushover.net/1/messages.json"
	icbcLoginUrl = "https://onlinebusiness.icbc.com/webdeas-ui/home"
)

func SendNotification(userKey, token, message string) error {
	_, err := http.PostForm(pushoverUrl, url.Values{
		"token":   {token},
		"user":    {userKey},
		"message": {"There is an available date: " + message + ", go: " + icbcLoginUrl + "."},
	})

	return err
}
