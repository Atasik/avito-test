package service

import (
	"segmenter/internal/domain"
	"segmenter/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpsertUserSegments(userID int, segmentsToAdd, segmentToDelete []domain.Segment) error {
	return s.userRepo.UpsertUserSegments(userID, segmentsToAdd, segmentToDelete)
}

func (s *UserService) GetSegments(userID int) ([]domain.Segment, error) {
	return s.userRepo.GetSegments(userID)
}

func (s *UserService) DeleteExpiredSegments() error {
	return s.userRepo.DeleteExpiredSegments()
}
