package order

import (
	"errors"
)

type Request struct {
	ID         string   `json:"id"`
	UserID     *string  `json:"user_id"`
	Products   []string `json:"products"`
	TotalPrice *float64 `json:"total_price"`
	Status     *string  `json:"status"`
}

func (s *Request) Validate() error {
	if s.UserID == nil {
		return errors.New("user_id: cannot be blank")
	}

	if s.Products == nil {
		return errors.New("products: cannot be blank")
	}

	if len(s.Products) == 0 {
		return errors.New("products: cannot be empty")
	}

	if s.TotalPrice == nil {
		return errors.New("total_price: cannot be blank")
	}

	if s.Status == nil {
		return errors.New("status: cannot be blank")
	} else if *s.Status != "new" && *s.Status != "processing" && *s.Status != "completed" {
		return errors.New("status: invalid value")
	}

	return nil
}

func (s *Request) IsEmpty(check string) error {
	if check == "update" {
		if s.UserID == nil && s.Products == nil &&
			s.TotalPrice == nil && s.Status == nil {
			return errors.New("data: cannot be blank")
		}

		if s.Products != nil && len(s.Products) == 0 {
			return errors.New("products: cannot be empty")
		}

		if s.Status != nil && *s.Status != "new" && *s.Status != "processing" && *s.Status != "completed" {
			return errors.New("status: invalid value")
		}
	}

	if check == "search" {
		if s.UserID == nil && s.Status == nil {
			return errors.New("invalid query: userId or status is required")
		}
	}

	return nil
}

type Response struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	Products   []string `json:"products"`
	TotalPrice float64  `json:"total_price"`
	Status     string   `json:"status"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:         data.ID,
		UserID:     *data.UserID,
		Products:   data.Products,
		TotalPrice: *data.TotalPrice,
		Status:     *data.Status,
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
