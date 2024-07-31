package epay

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strconv"
)

type PaymentResponse struct {
	ID             string  `json:"id"`
	AccountID      string  `json:"accountId"`
	Amount         float64 `json:"amount"`
	AmountBonus    float64 `json:"amountBonus"`
	Currency       string  `json:"currency"`
	Description    string  `json:"description"`
	Email          string  `json:"email"`
	InvoiceID      string  `json:"invoiceID"`
	Language       string  `json:"language"`
	Phone          string  `json:"phone"`
	Reference      string  `json:"reference"`
	IntReference   string  `json:"intReference"`
	Secure3D       *string `json:"secure3D"`
	Fingerprint    *string `json:"fingerprint"`
	CardID         string  `json:"cardID"`
	Fee            float64 `json:"fee"`
	ApprovalCode   string  `json:"approvalCode"`
	Code           int     `json:"code"`
	Status         string  `json:"status"`
	Secure3DStatus string  `json:"secure3DStatus"`
}

type PaymentReq struct {
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
	Name            string `json:"name"`
	Cryptogram      string `json:"cryptogram"`
	InvoiceID       string `json:"invoiceId"`
	InvoiceIDAlt    string `json:"invoiceIdAlt"`
	Description     string `json:"description"`
	AccountID       string `json:"accountId"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	CardSave        bool   `json:"cardSave"`
	Data            string `json:"data"`
	PostLink        string `json:"postLink"`
	FailurePostLink string `json:"failurePostLink"`
}

var (
	ClientID           = "test"
	ClientSecret       = "yF587AV9Ms94qN2QShFzVR3vFnWkhjbAK3sG"
	DefaultPaymentData = `{
 		"hpan":"4405639704015096",
		"expDate":"0125",
		"cvc":"815",
		"terminalId":
		"67e34d63-102f-4bd1-898e-370781d0074d"
	}`
)

func (c *Client) encryptData(data string) (string, error) {
	path, err := url.Parse(c.credentials.URL)
	if err != nil {
		return "", err
	}

	path = path.JoinPath("/public.rsa")

	resp, err := c.httpClient.Get(path.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(body)
	if block == nil {
		return "", errors.New("failed to decode pem block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (c *Client) Pay(ctx context.Context, token string, req *PaymentRequest) (dst PaymentResponse, err error) {
	u, err := url.Parse(c.credentials.URL)
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "payment/cryptopay")

	encryptedData, err := c.encryptData(DefaultPaymentData)
	amount, err := strconv.Atoi(req.Amount)
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return
	}
	if err != nil {
		return
	}
	paymentData := PaymentReq{
		Amount:          amount,
		Currency:        req.Currency,
		Name:            "JON JONSON",
		Cryptogram:      encryptedData,
		InvoiceID:       req.InvoiceID,
		InvoiceIDAlt:    "1",
		Description:     "test payment",
		AccountID:       "uuid000001",
		Email:           "jj@example.com",
		Phone:           "77777777777",
		CardSave:        true,
		Data:            "{\\\"statement\\\":{\\\"name\\\":\\\"Arman Ali\\\",\\\"invoiceID\\\":\\\"80000016\\\"}}",
		PostLink:        "https://testmerchant/order/1123",
		FailurePostLink: "https://testmerchant/order/1123/fail",
	}
	//data := url.Values{}
	//data.Set("amount", "100")
	//data.Set("currency", "KZT")
	//data.Set("name", "JON JONSON")
	//data.Set("cryptogram", encryptedData)
	//data.Set("invoiceId", invoiceId)
	//data.Set("invoiceIdAlt", "1")
	//data.Set("description", "test payment")
	//data.Set("accountId", "uuid000001")
	//data.Set("email", "jj@example.com")
	//data.Set("phone", "77777777777")
	//// data.Set("cardSave", "true")
	//// data.Set("data", "{\"statement\":{\"name\":\"Arman Ali\",\"invoiceID\":\"80000016\"}}")
	//// data.Set("postLink", "https://testmerchant/order/1123")
	//// data.Set("failurePostLink", "https://testmerchant/order/1123/fail")
	data, err := json.Marshal(paymentData)
	body := bytes.NewBuffer(data)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	err = c.request(ctx, "POST", u.String(), body, headers, &dst)

	return
}
