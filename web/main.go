package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fjacquet/selma-cli/internal/csvprocessor"
	"github.com/fjacquet/selma-cli/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Execute runs the root command.
func main() {
	logger.SetupLogger()

	// Initialize the Gin router
	r := gin.Default()

	// Serve static files (e.g., HTML, CSS, JS)
	r.Static("/static", "./static")

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	r.POST("/upload", uploadHandler)
	r.GET("/download/:filename", downloadHandler)

	// Start the server
	port := "8080"
	logrus.Infof("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func uploadHandler(c *gin.Context) {
	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("Failed to upload file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	// Save the uploaded file
	origFilePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, origFilePath); err != nil {
		logrus.Errorf("Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Process the file
	records, err := csvprocessor.ReadCSV(origFilePath)
	if err != nil {
		logrus.Errorf("Failed to read CSV: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV"})
		return
	}

	processedRecords := csvprocessor.ProcessRecords(records)

	// Save the processed records to a new CSV file
	outputFileName := fmt.Sprintf("processed_%d.csv", time.Now().Unix())
	outputFilePath := filepath.Join("downloads", outputFileName)
	if err := csvprocessor.WriteCSV(outputFilePath, processedRecords); err != nil {
		logrus.Errorf("Failed to write CSV: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV"})
		return
	}

	// Stream the processed CSV file back to the user
	c.Header("Content-Disposition", "attachment; filename=processed.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write the CSV headers
	headers := []string{"Date", "Description", "Bookkeeping No.", "Fund", "Amount", "Currency", "Number of Shares", "Stamp Duty Amount", "Investment"}
	if err := writer.Write(headers); err != nil {
		logrus.Errorf("Failed to write CSV headers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers"})
		return
	}

	// Write each record to the CSV
	for _, record := range processedRecords {
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
			logrus.Errorf("Failed to write CSV row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
			return
		}
	}

	logrus.Info("Processed file successfully streamed to client")

	logrus.Infof("File processed successfully: %s", outputFileName)
	c.JSON(http.StatusOK, gin.H{"message": "File processed successfully", "download": "/download/" + outputFileName})
}

func downloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("downloads", filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
