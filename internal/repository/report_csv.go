package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"segmenter/internal/domain"
	"time"
)

type ReportCSVRepository struct {
	ReportsDir string
}

func NewReportCSVRepository(reportsDir string) *ReportCSVRepository {
	return &ReportCSVRepository{
		ReportsDir: reportsDir,
	}
}

type ReportRepository interface {
	SaveReport(reportName string, history []domain.History) (string, error)
}

func (repo *ReportCSVRepository) SaveReport(reportName string, history []domain.History) (string, error) {
	reportName = fmt.Sprintf("%s.csv", reportName)
	filePath := filepath.Join(repo.ReportsDir, reportName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, operation := range history {
		row := []string{fmt.Sprintf("%d", operation.UserID), operation.Segment, operation.Operation, operation.CreatedAt.Format(time.DateTime)}
		err := writer.Write(row)
		if err != nil {
			return "", err
		}
	}

	return reportName, nil
}
