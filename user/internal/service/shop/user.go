package shop

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) ListUsers(ctx context.Context) (res []user.Response, err error) {
	data, err := s.userRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v\n", err)
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) CreateUser(ctx context.Context, req user.Request) (res user.Response, err error) {
	data := user.Entity{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Role:    req.Role,
	}

	data.ID, err = s.userRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("faled to create: %v\n", err)
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) GetUser(ctx context.Context, id string) (res user.Response, err error) {
	data, err := s.userRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v\n", err)
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateUser(ctx context.Context, id string, req user.Request) (err error) {
	data := user.Entity{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Role:    req.Role,
	}

	err = s.userRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteUser(ctx context.Context, id string) (err error) {
	err = s.userRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchUser(ctx context.Context, req user.Request) (res []user.Response, err error) {
	searchData := user.Entity{
		Name:  req.Name,
		Email: req.Email,
	}
	data, err := s.userRepository.Search(ctx, searchData)
	if err != nil {
		fmt.Printf("failed to search users: %v\n", err)
		return
	}

	res = user.ParseFromEntities(data)

	return
}
