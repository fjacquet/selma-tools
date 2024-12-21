package models

// Record represents a financial transaction or investment record.
type Record struct {
	Date            string  // Date of the transaction in YYYY-MM-DD format.
	Description     string  // Brief description of the transaction.
	BookkeepingNo   string  // Unique identifier for bookkeeping purposes.
	Fund            string  // The fund or account associated with the transaction.
	Amount          float64 // The monetary value of the transaction.
	Currency        string  // Currency in which the Amount is denominated, e.g., USD, EUR.
	NumberOfShares  string  // Quantity of shares involved, if applicable.
	StampDutyAmount float64 // Stamp duty or tax amount applied to the transaction.
	Investment      string  // Details about the related investment, such as stock or bond name.
}
