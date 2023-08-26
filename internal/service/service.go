package service

import "segmenter/internal/repository"

type Service struct {
	Segment
	User
	History
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Segment: NewSegmentService(repos),
		User:    NewUserService(repos),
		History: NewHistoryService(repos),
	}
}
