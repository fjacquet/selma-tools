package csvprocessor

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/fjacquet/selma-tools/models"
	"github.com/sirupsen/logrus"
)

// ReadCSV reads and parses a CSV file into a slice of Record objects.
func ReadCSV(filePath string) ([]models.Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filePath": filePath}).Errorf("Failed to open CSV file: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		logrus.WithFields(logrus.Fields{"filePath": filePath}).Errorf("Failed to read CSV headers: %v", err)
		return nil, err
	}

	if len(headers) < 7 {
		logrus.WithFields(logrus.Fields{"filePath": filePath, "headers": headers}).Error("Unexpected CSV header length")
		return nil, errors.New("unexpected CSV header length")
	}

	rows, err := reader.ReadAll()
	if err != nil {
		logrus.WithFields(logrus.Fields{"filePath": filePath}).Errorf("Failed to read CSV rows: %v", err)
		return nil, err
	}

	var records []models.Record
	for _, row := range rows {
		record, err := parseCSVRow(row)
		if err != nil {
			logrus.WithFields(logrus.Fields{"row": row}).Warnf("Failed to parse CSV row: %v", err)
			continue
		}
		records = append(records, record)
	}

	logrus.WithFields(logrus.Fields{"filePath": filePath, "count": len(records)}).Info("Successfully read CSV records")
	return records, nil
}

// parseCSVRow parses a single CSV row into a Record object.
func parseCSVRow(row []string) (models.Record, error) {
	amount, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{
		Date:            row[0],
		Description:     row[1],
		BookkeepingNo:   row[2],
		Fund:            row[3],
		Amount:          amount,
		Currency:        row[5],
		NumberOfShares:  row[6],
		StampDutyAmount: 0,
		Investment:      "",
	}, nil
}
