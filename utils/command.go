package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCommandToFile runs a command and redirects its output to a specified file.
func RunCommandToFile(command string, args []string, outputFile string) error {
	cmd := exec.Command(command, args...)
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer output.Close()

	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command %s: %w", command, err)
	}
	return nil
}

// RunCommand runs a command and outputs the results to the terminal.
func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunShellCommand runs a shell command and outputs its results.
func RunShellCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}
