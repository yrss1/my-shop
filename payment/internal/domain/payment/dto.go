package payment

import (
	"errors"
)

type Request struct {
	ID      string  `json:"id"`
	UserID  *string `json:"user_id"`
	OrderID *string `json:"order_id"`
	Amount  *string `json:"amount"`
	Status  *string `json:"status"`
}

func (s *Request) Validate() error {
	if s.UserID == nil {
		return errors.New("user_id: cannot be blank")
	}

	if s.OrderID == nil {
		return errors.New("order_id: cannot be blank")
	}

	if s.Amount == nil {
		return errors.New("amount: cannot be blank")
	}

	//if s.Status == nil {
	//	return errors.New("status: cannot be blank")
	//} else if *s.Status != "new" && *s.Status != "processing" && *s.Status != "completed" {
	//	return errors.New("status: invalid value")
	//}

	return nil
}

func (s *Request) IsEmpty(check string) error {
	if check == "update" {
		if s.UserID == nil && s.OrderID == nil &&
			s.Amount == nil && s.Status == nil {
			return errors.New("data: cannot be blank")
		}
		if s.Status != nil && *s.Status != "successful" && *s.Status != "unsuccessful" && *s.Status != "pending" {
			return errors.New("status: invalid value")
		}
	}

	if check == "search" {
		if s.OrderID == nil && s.Status == nil && s.UserID == nil {
			return errors.New("invalid query")
		}
	}

	return nil
}

type Response struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
	Amount  string `json:"amount"`
	Status  string `json:"status"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:      data.ID,
		UserID:  *data.UserID,
		OrderID: *data.OrderID,
		Amount:  *data.Amount,
		Status:  *data.Status,
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
