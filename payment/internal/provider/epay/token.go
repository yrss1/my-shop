package epay

import (
	"bytes"
	"context"
	"github.com/shopspring/decimal"
	"net/url"
	"time"
)

type TokenResponse struct {
	Scope        string          `json:"scope"`
	ExpiresIn    decimal.Decimal `json:"expires_in"`
	TokenType    string          `json:"token_type"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
}

func (c *Client) initGlobalTokenRefresher() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.credentials.GlobalToken, err = c.GetPaymentToken(ctx, nil)
	if err != nil {
		return
	}
	ticker := time.NewTicker(time.Duration(c.credentials.GlobalToken.ExpiresIn.IntPart()-60) * time.Second)

	go func() {
		for {
			<-ticker.C
			c.credentials.GlobalToken, err = c.GetPaymentToken(ctx, nil)
		}
	}()
	return
}

func (c *Client) GetPaymentToken(ctx context.Context, src *PaymentRequest) (dst TokenResponse, err error) {
	parsedURL, err := url.Parse(c.credentials.OAuthURL)
	if err != nil {
		return
	}
	parsedURL.Path = "/epay2/oauth2/token"

	data := url.Values{}
	data.Set("client_id", c.credentials.Login)
	data.Set("client_secret", c.credentials.Password)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "webapi usermanagement email_send verification statement statistics payment")
	if src != nil {
		data.Set("invoiceID", src.InvoiceID)
		data.Set("amount", src.Amount)
		data.Set("currency", src.Currency)
		data.Set("terminal", src.TerminalID)
	}

	body := bytes.NewBufferString(data.Encode())

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err = c.request(ctx, "POST", parsedURL.String(), body, headers, &dst)

	return
}
