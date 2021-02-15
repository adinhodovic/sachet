package clickatell

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for the Clickatell provider
type ClickatellConfig struct {
	Token string `yaml:"token"`
}

// Clickatell contains the necessary values for the Clickatell provider
type Clickatell struct {
	ClickatellConfig
}

type Payload struct {
	Text string   `json:"text"`
	To   []string `json:"to"`
	From string   `json:"from,omitempty"`
}

// NewClickatell creates and returns a new Clickatell struct
func NewClickatell(config ClickatellConfig) *Clickatell {
	return &Clickatell{config}
}

// Send messages to receipients
func (c *Clickatell) Send(message sachet.Message) error {
	payloadBuf := new(bytes.Buffer)
	payload := &Payload{Text: message.Text, To: message.To, From: message.From}
	jsonErr := json.NewEncoder(payloadBuf).Encode(payload)
	if jsonErr != nil {
		return jsonErr
	}

	req, err := http.NewRequest("POST", "https://api.clickatell.com/rest/message", payloadBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer: "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, resErr := client.Do(req)

	if resErr != nil {
		return resErr
	}
	defer res.Body.Close()

	return nil
}
