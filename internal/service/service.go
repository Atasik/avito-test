package service

import (
	"segmenter/internal/domain"
	"segmenter/internal/repository"
	"time"
)

type Service struct {
	Segment
	User
	History
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Segment: NewSegmentService(repos.SegmentRepo),
		User:    NewUserService(repos.UserRepo),
		History: NewHistoryService(repos.HistoryRepo, repos.ReportRepository),
	}
}

type History interface {
	CreateReport(period time.Time, userID int) (string, error)
}

type User interface {
	UpsertUserSegments(userID int, segmentsToAdd, segmentToDelete []domain.Segment) error
	GetSegments(userID int) ([]domain.Segment, error)
	DeleteExpiredSegments() error
}

type Segment interface {
	Create(seg domain.Segment) (int, error)
	Delete(seg domain.Segment) error
}
