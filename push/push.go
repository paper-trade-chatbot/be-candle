package push

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/paper-trade-chatbot/be-candle/config"
	"github.com/paper-trade-chatbot/be-candle/logging"
)

// Message ...
type Message struct {
	To           string            `json:"to"`
	Data         map[string]string `json:"data"`
	Notification FCMNotification   `json:"notification"`
}

// FCMNotification is the notification message body
type FCMNotification struct {
	BodyLocKey  string   `json:"body_loc_key"`
	BodyLocArgs []string `json:"body_loc_args"`
	Badge       int      `json:"badge"`
}

const (
	fcmURL = "https://fcm.googleapis.com/fcm/send"
)

var (
	fcmkey = config.GetString("FCM_KEY")
)

// Send ...
func Send(msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer([]byte(b))

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", "key="+fcmkey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logging.Debug(string(b))

	return nil
}
