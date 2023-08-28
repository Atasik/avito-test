package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"segmenter/internal/repository"
	"time"
)

type History interface {
	CreateReport(period time.Time, userID int) (string, error)
}

type HistoryService struct {
	historyRepo repository.HistoryRepo
}

func NewHistoryService(historyRepo repository.HistoryRepo) *HistoryService {
	return &HistoryService{historyRepo: historyRepo}
}

func (s *HistoryService) CreateReport(period time.Time, userID int) (string, error) {
	history, err := s.historyRepo.GetHistoryForPeriod(period, userID)
	if err != nil {
		return "", err
	}

	// TODO: refactor, remove hardcode
	reportsPath := "./reports/"
	reportName := fmt.Sprintf("%d.csv", time.Now().Unix())

	file, err := os.Create(fmt.Sprintf("%s%s", reportsPath, reportName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, operation := range history {
		row := []string{fmt.Sprintf("%d", operation.UserID), operation.Segment, operation.Operation, time.Time(operation.CreatedAt).Format(time.DateTime)}
		err := writer.Write(row)
		if err != nil {
			return "", err
		}
	}

	return reportName, nil
}
