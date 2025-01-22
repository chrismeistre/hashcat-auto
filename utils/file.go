package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

// WriteToFile writes a slice of strings to a file, one string per line.
func WriteToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer for file %s: %w", filename, err)
	}
	return nil
}

// AppendToFile appends a slice of strings to a file, one string per line.
// If the file does not exist, it creates a new file.
func AppendToFile(filename string, lines []string) error {
	// Open the file in append mode, create it if it does not exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer for file %s: %w", filename, err)
	}
	return nil
}

// ValidateFileExists checks if a file exists at the given path.
func ValidateFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filepath)
	}
	return nil
}

func ValidateFiles(files []string) error {
	for _, file := range files {
		if err := ValidateFileExists(file); err != nil {
			return fmt.Errorf("file validation failed for %s: %w", file, err)
		}
	}
	color.Green("All provided files exist.")
	return nil
}

func CountLines(filename string) (int, error) {
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		return len(lines), nil
	}
	return 0, nil
}

func CompareCrackedFiles(currentFile, previousFile string) (int, error) {
	currentCount, err := CountLines(currentFile)
	if err != nil {
		return 0, err
	}

	previousCount := 0
	if _, err := os.Stat(previousFile); err == nil { // Previous file exists
		previousCount, err = CountLines(previousFile)
		if err != nil {
			return 0, err
		}
	}

	return currentCount - previousCount, nil
}
