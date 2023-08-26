package service

import (
	"segmenter/internal/repository"
)

type User interface {
	UpsertUser(id int, segmentsToAdd, segmentToDelete []repository.Segment) error
	GetUserSegments(id int) ([]repository.Segment, error)
}

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpsertUser(id int, segmentsToAdd, segmentToDelete []repository.Segment) error {
	return s.userRepo.UpsertUser(id, segmentsToAdd, segmentToDelete)
}

func (s *UserService) GetUserSegments(id int) ([]repository.Segment, error) {
	return s.userRepo.GetUserSegments(id)
}
