package epay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Credentials struct {
	URL      string
	Login    string
	Password string

	OAuthURL       string
	PaymentPageURL string
	GlobalToken    TokenResponse
}

type Client struct {
	httpClient  *http.Client
	credentials Credentials
}

type PaymentRequest struct {
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	InvoiceID  string `json:"invoice_id"`
	TerminalID string `json:"terminal_id"`
}

func New(credentials Credentials) (client *Client, err error) {
	httpClient := http.DefaultClient
	httpClient.Timeout = 30 * time.Second

	client = &Client{
		httpClient:  httpClient,
		credentials: credentials,
	}
	err = client.initGlobalTokenRefresher()

	return
}

func (c *Client) request(ctx context.Context, method, url string, body *bytes.Buffer, headers map[string]string, dst interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}
