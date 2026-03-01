package service

import (
	"golf-store/be-mono/internal/modules/user/repository"
	db "golf-store/be-mono/internal/platform/db"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserWithPasswordByID(id string) (*db.UserEntity, error) {
	user, err := s.repo.GetUserWithPasswordByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(id string) (*db.UserEntity, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
