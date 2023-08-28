package service

import (
	"segmenter/internal/domain"
	"segmenter/internal/repository"
)

type Segment interface {
	CreateSegment(seg domain.Segment) (int, error)
	DeleteSegment(seg domain.Segment) error
}

type SegmentService struct {
	segmentRepo repository.SegmentRepo
}

func NewSegmentService(segmentRepo repository.SegmentRepo) *SegmentService {
	return &SegmentService{segmentRepo: segmentRepo}
}

func (s *SegmentService) CreateSegment(seg domain.Segment) (int, error) {
	return s.segmentRepo.CreateSegment(seg)
}

func (s *SegmentService) DeleteSegment(seg domain.Segment) error {
	return s.segmentRepo.DeleteSegment(seg)
}
