package service

import "segmenter/internal/repository"

type Segment interface {
	CreateSegment(slug repository.Segment) (int, error)
	DeleteSegment(slug repository.Segment) error
}

type SegmentService struct {
	segmentRepo repository.SegmentRepo
}

func NewSegmentService(segmentRepo repository.SegmentRepo) *SegmentService {
	return &SegmentService{segmentRepo: segmentRepo}
}

func (s *SegmentService) CreateSegment(slug repository.Segment) (int, error) {
	return s.segmentRepo.CreateSegment(slug)
}

func (s *SegmentService) DeleteSegment(slug repository.Segment) error {
	return s.segmentRepo.DeleteSegment(slug)
}
