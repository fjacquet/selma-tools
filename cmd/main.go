package main

import (
	"github.com/fjacquet/selma-tools/internal/csvprocessor"
	"github.com/fjacquet/selma-tools/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var origFile string
var cleanFile string

// Execute runs the root command.
func main() {
	logger.SetupLogger()

	rootCmd := createRootCmd(&origFile, &cleanFile)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

// createRootCmd creates the root command for the CLI application.
func createRootCmd(origFile, cleanFile *string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "selma-cli",
		Short: "Process CSV files to add investment and stamp duty details",
		Run: func(cmd *cobra.Command, args []string) {
			processCSVFiles(*origFile, *cleanFile)
		},
	}

	rootCmd.Flags().StringVarP(origFile, "input", "i", "", "Path to the original CSV file (required)")
	rootCmd.Flags().StringVarP(cleanFile, "output", "o", "", "Path to the output CSV file (required)")

	if err := rootCmd.MarkFlagRequired("input"); err != nil {
		logrus.Fatalf("Failed to mark input flag as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("output"); err != nil {
		logrus.Fatalf("Failed to mark output flag as required: %v", err)
	}

	return rootCmd
}

func processCSVFiles(origFile, cleanFile string) {
	records, err := csvprocessor.ReadCSV(origFile)
	if err != nil {
		logrus.Errorf("Error reading CSV: %v", err)
		return
	}

	processedRecords := csvprocessor.ProcessRecords(records)

	if err := csvprocessor.WriteCSV(cleanFile, processedRecords); err != nil {
		logrus.Errorf("Error writing CSV: %v", err)
	}
}
