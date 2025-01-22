package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// ExtractUsernames extracts usernames from a hashlist file, handling DOMAIN\\username formats.
func ExtractUsernames(hashlist string) ([]string, error) {
	file, err := os.Open(hashlist)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", hashlist, err)
	}
	defer file.Close()

	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 0 {
			username := parts[0]
			if strings.Contains(username, "\\") {
				usernameParts := strings.Split(username, "\\")
				if len(usernameParts) > 1 {
					username = usernameParts[1]
				}
			}
			usernames = append(usernames, username)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file %s: %w", hashlist, err)
	}
	return usernames, nil
}

// ExtractPasswords parses a file and extracts unique, non-empty passwords from `user:password` lines.
func ExtractPasswords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	passwordSet := make(map[string]struct{}) // To ensure uniqueness
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			password := strings.TrimSpace(parts[len(parts)-1]) // Extract the last part as the password
			if password != "" {                                // Skip empty passwords
				passwordSet[password] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file %s: %w", filename, err)
	}

	// Convert the set to a slice
	passwords := make([]string, 0, len(passwordSet))
	for password := range passwordSet {
		passwords = append(passwords, password)
	}

	return passwords, nil
}

// CleanWordlist removes all non-ASCII characters from a wordlist file and adds a timestamp to the output filename.
func CleanWordlist(inputFile, outputDir, timestamp string) (string, error) {
	input, err := os.Open(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to open input file %s: %w", inputFile, err)
	}
	defer input.Close()

	outputFile := fmt.Sprintf("%s/cleaned_wordlist_%s.txt", outputDir, timestamp)
	output, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	for scanner.Scan() {
		line := scanner.Text()
		cleaned := removeNonASCII(line)
		if cleaned != "" {
			if _, err := writer.WriteString(cleaned + "\n"); err != nil {
				return "", fmt.Errorf("failed to write to output file %s: %w", outputFile, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan input file %s: %w", inputFile, err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush output file %s: %w", outputFile, err)
	}

	return outputFile, nil
}

// removeNonASCII removes all non-ASCII characters from a string.
func removeNonASCII(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if r <= unicode.MaxASCII {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
