package epay

import (
	"context"
	"net/url"
	"time"
)

type TransactionResponse struct {
	ID                string    `json:"id"`
	CreatedDate       time.Time `json:"createdDate"`
	InvoiceID         string    `json:"invoiceID"`
	Amount            int       `json:"amount"`
	AmountBonus       int       `json:"amountBonus"`
	OrgAmount         int       `json:"orgAmount"`
	ApprovalCode      string    `json:"approvalCode"`
	PayoutAmount      int       `json:"payoutAmount"`
	Currency          string    `json:"currency"`
	Terminal          string    `json:"terminal"`
	AccountID         string    `json:"accountID"`
	Description       string    `json:"description"`
	Data              string    `json:"data"`
	Language          string    `json:"language"`
	CardMask          string    `json:"cardMask"`
	CardType          string    `json:"cardType"`
	Issuer            string    `json:"issuer"`
	Reference         string    `json:"reference"`
	Reason            string    `json:"reason"`
	ReasonCode        string    `json:"reasonCode"`
	IntReference      string    `json:"intReference"`
	Secure            bool      `json:"secure"`
	StatusID          string    `json:"statusID"`
	StatusName        string    `json:"statusName"`
	StatusDescription string    `json:"statusDescription"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	CardID            string    `json:"cardID"`
	XlsRRN            string    `json:"xlsRRN"`
	IP                string    `json:"ip"`
	IPCountry         string    `json:"ipCountry"`
	IPCity            string    `json:"ipCity"`
	IPRegion          string    `json:"ipRegion"`
	IPDistrict        string    `json:"ipDistrict"`
	IPLatitude        float64   `json:"ipLatitude"`
	IPLongitude       float64   `json:"ipLongitude"`
}

type StatusResponse struct {
	ResultCode    string              `json:"resultCode"`
	ResultMessage string              `json:"resultMessage"`
	Transaction   TransactionResponse `json:"transaction"`
}

func (c *Client) GetStatus(ctx context.Context, token string, invoiceID string) (dst StatusResponse, err error) {
	path, err := url.Parse(c.credentials.URL)
	if err != nil {
		return
	}

	path = path.JoinPath("/check-status/payment/transaction/", invoiceID)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	if err = c.request(ctx, "GET", path.String(), nil, headers, &dst); err != nil {
		return
	}

	return
}

//func (c *Client) GetPaymentStatus(ctx context.Context, token string, invoiceID string) (dst StatusResponse, err error) {
//	path, err := url.Parse(c.credentials.URL)
//	if err != nil {
//		return
//	}
//
//	path = path.JoinPath("/check-status/payment/transaction/", invoiceID)
//
//	headers := map[string]string{
//		"Content-Type":  "application/json",
//		"Authorization": fmt.Sprintf("Bearer %s", token),
//	}
//	if err = c.request(ctx, "POST", path.String(), nil, headers, &dst); err != nil {
//		return
//	}
//	return
//}
