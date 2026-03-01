package service

import (
	"golf-store/be-mono/internal/modules/user/repository"

	entities "golf-store/be-mono/internal/platform/entities"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserByID(id string) (*entities.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
