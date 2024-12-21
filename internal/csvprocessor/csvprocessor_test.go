package csvprocessor

import (
	"testing"

	"github.com/fjacquet/selma-tools/models"
	"github.com/stretchr/testify/assert"
)

func TestProcessRecords(t *testing.T) {
	records := []models.Record{
		{Date: "2023-10-01", Description: "trade", Amount: -1000.00},
		{Date: "2023-10-01", Description: "stamp_duty", Amount: 5.00},
		{Date: "2023-10-02", Description: "dividend", Amount: 50.00},
		{Date: "2023-10-03", Description: "cash_transfer", Amount: 200.00},
		{Date: "2023-10-04", Description: "selma_fee", Amount: -10.00},
	}

	expected := []models.Record{
		{Date: "2023-10-01", Description: "trade", Amount: -1000.00, Investment: "Buy", StampDutyAmount: 5.00},
		{Date: "2023-10-02", Description: "dividend", Amount: 50.00, Investment: "Dividend"},
		{Date: "2023-10-03", Description: "cash_transfer", Amount: 200.00, Investment: "Income"},
		{Date: "2023-10-04", Description: "selma_fee", Amount: -10.00, Investment: "Expense"},
	}

	result := ProcessRecords(records)
	assert.Equal(t, expected, result)
}

func TestCategorizeRecord(t *testing.T) {
	tests := []struct {
		name     string
		record   models.Record
		expected models.Record
	}{
		{"Dividend", models.Record{Description: "dividend"}, models.Record{Description: "dividend", Investment: "Dividend"}},
		{"Cash Transfer", models.Record{Description: "cash_transfer"}, models.Record{Description: "cash_transfer", Investment: "Income"}},
		{"Selma Fee", models.Record{Description: "selma_fee"}, models.Record{Description: "selma_fee", Investment: "Expense"}},
		{"Trade Buy", models.Record{Description: "trade", Amount: -1000.00}, models.Record{Description: "trade", Amount: -1000.00, Investment: "Buy"}},
		{"Trade Sell", models.Record{Description: "trade", Amount: 1000.00}, models.Record{Description: "trade", Amount: 1000.00, Investment: "Sell"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := categorizeRecord(tt.record)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleStampDuty(t *testing.T) {
	tests := []struct {
		name     string
		record   models.Record
		previous *models.Record
		next     *models.Record
		expected models.Record
	}{
		{"Next Stamp Duty", models.Record{Description: "trade"}, nil, &models.Record{Description: "stamp_duty", Amount: 5.00}, models.Record{Description: "trade", StampDutyAmount: 5.00}},
		{"Previous Stamp Duty", models.Record{Description: "trade"}, &models.Record{Description: "stamp_duty", Amount: 5.00}, nil, models.Record{Description: "trade", StampDutyAmount: 5.00}},
		{"No Stamp Duty", models.Record{Description: "trade"}, nil, nil, models.Record{Description: "trade"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleStampDuty(tt.record, tt.previous, tt.next)
			assert.Equal(t, tt.expected, result)
		})
	}
}
