package service

import (
	"segmenter/internal/domain"
	"segmenter/internal/repository"
)

type User interface {
	UpsertUser(userID int, segmentsToAdd, segmentToDelete []domain.Segment) error
	GetUserSegments(userID int) ([]domain.Segment, error)
	DeleteExpiredSegments() error
}

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpsertUser(userID int, segmentsToAdd, segmentToDelete []domain.Segment) error {
	return s.userRepo.UpsertUser(userID, segmentsToAdd, segmentToDelete)
}

func (s *UserService) GetUserSegments(userID int) ([]domain.Segment, error) {
	return s.userRepo.GetUserSegments(userID)
}

func (s *UserService) DeleteExpiredSegments() error {
	return s.userRepo.DeleteExpiredSegments()
}
