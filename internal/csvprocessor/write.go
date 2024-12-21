package csvprocessor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/fjacquet/selma-tools/models"
	"github.com/sirupsen/logrus"
)

// WriteCSV writes a slice of Record structs to a CSV file.
func WriteCSV(filePath string, records []models.Record) error {
	file, err := os.Create(filePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filePath": filePath}).Errorf("Failed to create CSV file: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Date", "Description", "Bookkeeping No.", "Fund", "Amount", "Currency", "Number of Shares", "Stamp Duty Amount", "Investment"}
	if err := writer.Write(headers); err != nil {
		logrus.WithFields(logrus.Fields{"filePath": filePath}).Errorf("Failed to write CSV headers: %v", err)
		return err
	}

	for _, record := range records {
		row := []string{
			record.Date,
			record.Description,
			record.BookkeepingNo,
			record.Fund,
			fmt.Sprintf("%.2f", record.Amount),
			record.Currency,
			record.NumberOfShares,
			fmt.Sprintf("%.2f", record.StampDutyAmount),
			record.Investment,
		}
		if err := writer.Write(row); err != nil {
			logrus.WithFields(logrus.Fields{"filePath": filePath, "count": len(records)}).Info("Successfully wrote CSV records")
			return err
		}
	}

	return nil
}
