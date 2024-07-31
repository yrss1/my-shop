package epayment

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Service) GetToken(ctx context.Context, req *epay.PaymentRequest) (token string, err error) {
	dest, err := s.epayClient.GetPaymentToken(ctx, req)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		//logger.Error("failed to get token", zap.Error(err))
		return
	}

	token = dest.AccessToken

	return
}

func (s *Service) GetStatus(ctx context.Context, token string, invoiceID string) (dest epay.StatusResponse, err error) {
	dest, err = s.epayClient.GetStatus(ctx, token, invoiceID)
	if err != nil {
		//logger.Error("failed to get status", zap.Error(err))
		return
	}

	return
}

func (s *Service) Pay(ctx context.Context, token string, req *epay.PaymentRequest) (dest epay.PaymentResponse, err error) {
	dest, err = s.epayClient.Pay(ctx, token, req)
	if err != nil {
		return
	}

	return
}
