package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// Capture the output
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer func() {
		logrus.SetOutput(os.Stderr)
	}()

	// Set up the command with test flags
	rootCmd := createRootCmd(&origFile, &cleanFile)
	rootCmd.SetArgs([]string{"--input", "../testdata/input.csv", "--output", "../testdata/output.csv"})

	// Execute the command
	err := rootCmd.Execute()
	assert.NoError(t, err)

	// Check the output
	output := buf.String()
	assert.Contains(t, output, "Processed records")
}

func TestMissingInputFlag(t *testing.T) {
	// Capture the output
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer func() {
		logrus.SetOutput(os.Stderr)
	}()

	// Redirect os.Stderr to capture the error output
	stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Set up the command without the required input flag
	rootCmd := createRootCmd(&origFile, &cleanFile)
	rootCmd.SetArgs([]string{"--output", "testdata/output.csv"})

	// Execute the command
	err := rootCmd.Execute()
	assert.Error(t, err)

	// Restore os.Stderr and close the pipe
	w.Close()
	os.Stderr = stderr

	// Read the captured error output
	var errBuf bytes.Buffer
	_, _ = errBuf.ReadFrom(r)

	// Check the output
	output := errBuf.String()
	assert.Contains(t, output, "required flag(s) \"input\" not set")
}

func TestMissingOutputFlag(t *testing.T) {
	// Capture the output
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer func() {
		logrus.SetOutput(os.Stderr)
	}()

	// Redirect os.Stderr to capture the error output
	stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Set up the command without the required output flag
	rootCmd := createRootCmd(&origFile, &cleanFile)
	rootCmd.SetArgs([]string{"--input", "testdata/input.csv"})

	// Execute the command
	err := rootCmd.Execute()
	assert.Error(t, err)

	// Restore os.Stderr and close the pipe
	w.Close()
	os.Stderr = stderr

	// Read the captured error output
	var errBuf bytes.Buffer
	_, _ = errBuf.ReadFrom(r)

	// Check the output
	output := errBuf.String()
	assert.Contains(t, output, "required flag(s) \"output\" not set")
}
