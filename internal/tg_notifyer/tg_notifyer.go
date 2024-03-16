package tgnotifyer

import (
	"fmt"
	"net/http"
	"net/url"
)

type TgNotifyer struct {
	token      string
	chatID     string
	apiURL     string
	baseApiURL string
}

type Option func(*TgNotifyer)

func WithToken(token string) Option {
	return func(w *TgNotifyer) {
		w.token = token
	}
}

func WithChatID(chatID string) Option {
	return func(w *TgNotifyer) {
		w.chatID = chatID
	}
}

func New(opts ...Option) *TgNotifyer {
	notifyer := &TgNotifyer{
		token:      "",
		chatID:     "",
		apiURL:     "",
		baseApiURL: "https://api.telegram.org/bot%s/sendMessage?chat_id=%s",
	}

	for _, opt := range opts {
		opt(notifyer)
	}
	if notifyer.token == "" || notifyer.chatID == "" {
		return notifyer
	}
	notifyer.apiURL = fmt.Sprintf(notifyer.baseApiURL, notifyer.token, notifyer.chatID)
	return notifyer
}

// Send message to Telegram chat via bot.
//
// Returns error if something went wrong.
func (t *TgNotifyer) Notify(message string) error {
	// Check if token and chat ID are set.
	if t.apiURL == "" {
		return fmt.Errorf("API URL is not set")
	}
	// Construct API URL.
	apiURL := fmt.Sprintf("%s&text=%s", t.apiURL, url.QueryEscape(message))

	// Send request.
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// No error, yay!
	return nil
}
