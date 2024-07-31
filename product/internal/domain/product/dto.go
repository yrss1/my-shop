package product

import (
	"errors"
)

type Request struct {
	ID          string   `json:"id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
	Quantity    *int     `json:"quantity"`
}

func (s *Request) Validate() error {
	if s.Name == nil {
		return errors.New("name: cannot be blank")
	}

	if s.Price == nil {
		return errors.New("price: cannot be blank")
	}

	if s.Quantity == nil {
		return errors.New("quantity: cannot be blank")
	}

	return nil
}

func (s *Request) IsEmpty() error {
	if s.Name == nil && s.Description == nil &&
		s.Price == nil && s.Category == nil &&
		s.Quantity == nil {
		return errors.New("data cannot be blank")
	}
	return nil
}

type Response struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Quantity    int     `json:"quantity"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		Name:        *data.Name,
		Description: *data.Description,
		Price:       *data.Price,
		Category:    *data.Category,
		Quantity:    *data.Quantity,
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
