package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"segmenter/internal/repository"
	"time"
)

type History interface {
	CreateReport(start, end time.Time, userId int) (string, error)
}

type HistoryService struct {
	historyRepo repository.HistoryRepo
}

func NewHistoryService(historyRepo repository.HistoryRepo) *HistoryService {
	return &HistoryService{historyRepo: historyRepo}
}

func (s *HistoryService) CreateReport(start, end time.Time, userId int) (string, error) {
	history, err := s.historyRepo.GetHistoryForPeriod(start, end, userId)
	if err != nil {
		return "", err
	}

	fmt.Println(history)

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
		row := []string{fmt.Sprintf("%d", operation.UserId), operation.Segment, operation.Operation, operation.CreatedAt.Format(time.DateTime)}
		writer.Write(row)
	}

	return reportName, nil
}
