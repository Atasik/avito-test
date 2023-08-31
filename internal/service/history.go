package service

import (
	"fmt"
	"segmenter/internal/repository"
	"time"
)

type HistoryService struct {
	historyRepo repository.HistoryRepo
	reportRepo  repository.ReportRepository
}

func NewHistoryService(historyRepo repository.HistoryRepo, reportRepo repository.ReportRepository) *HistoryService {
	return &HistoryService{historyRepo: historyRepo, reportRepo: reportRepo}
}

func (s *HistoryService) CreateReport(period time.Time, userID int) (string, error) {
	history, err := s.historyRepo.GetHistoryForPeriod(period, userID)
	if err != nil {
		return "", err
	}

	reportName := fmt.Sprintf("%d", time.Now().Unix())
	reportName, err = s.reportRepo.SaveReport(reportName, history)
	if err != nil {
		return "", err
	}

	return reportName, nil
}
