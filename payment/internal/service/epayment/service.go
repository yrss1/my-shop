package epayment

type Configuration func(s *Service) error

type Service struct {
	paymentRepository payment.Repository
	epayClient        *epay.Client
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPaymentRepository(paymentRepository payment.Repository) Configuration {
	return func(s *Service) error {
		s.paymentRepository = paymentRepository
		return nil
	}
}

func WithEpayClient(epayClient *epay.Client) Configuration {
	return func(s *Service) error {
		s.epayClient = epayClient
		return nil
	}
}
