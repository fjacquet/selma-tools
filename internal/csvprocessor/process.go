package csvprocessor

import (
	"github.com/fjacquet/selma-tools/models"
	"github.com/sirupsen/logrus"
)

// ProcessRecords processes a slice of Record objects.
func ProcessRecords(records []models.Record) []models.Record {
	var newRecords []models.Record

	for i, record := range records {
		var previous, next *models.Record
		if i > 0 {
			previous = &records[i-1]
		}
		if i < len(records)-1 {
			next = &records[i+1]
		}

		record = categorizeRecord(record)
		if record.Description == "trade" {
			record = handleStampDuty(record, previous, next)
		}

		if record.Description != "stamp_duty" {
			newRecords = append(newRecords, record)
		}
	}

	logrus.WithFields(logrus.Fields{"inputCount": len(records), "outputCount": len(newRecords)}).Info("Processed records")
	return newRecords
}

// categorizeRecord categorizes the investment type of a record based on its description.
func categorizeRecord(record models.Record) models.Record {
	switch record.Description {
	case "dividend":
		record.Investment = "Dividend"
	case "cash_transfer":
		record.Investment = "Income"
	case "selma_fee":
		record.Investment = "Expense"
	case "trade":
		if record.Amount < 0 {
			record.Investment = "Buy"
		} else {
			record.Investment = "Sell"
		}
	}
	return record
}

// handleStampDuty associates stamp duty amounts with trade records.
func handleStampDuty(record models.Record, previous, next *models.Record) models.Record {
	if next != nil && next.Description == "stamp_duty" {
		record.StampDutyAmount = next.Amount
	} else if previous != nil && previous.Description == "stamp_duty" {
		record.StampDutyAmount = previous.Amount
	}
	return record
}
