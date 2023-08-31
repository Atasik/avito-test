package service

import (
	"segmenter/internal/domain"
	"segmenter/internal/repository"
)

type SegmentService struct {
	segmentRepo repository.SegmentRepo
}

func NewSegmentService(segmentRepo repository.SegmentRepo) *SegmentService {
	return &SegmentService{segmentRepo: segmentRepo}
}

func (s *SegmentService) Create(seg domain.Segment) (int, error) {
	return s.segmentRepo.Create(seg)
}

func (s *SegmentService) Delete(seg domain.Segment) error {
	return s.segmentRepo.Delete(seg)
}
