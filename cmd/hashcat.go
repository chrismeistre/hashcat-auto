package cmd

import (
	"fmt"
	"hashcat-auto/config"
	"hashcat-auto/utils"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

func validateEnvironment(hashcatPath string) error {
	color.Yellow("Validating environment...")

	// Check if Hashcat is installed
	if _, err := exec.LookPath(hashcatPath); err != nil {
		return fmt.Errorf("Hashcat not found at %s: %w", hashcatPath, err)
	}
	color.Green("Hashcat is installed: %s", hashcatPath)

	// Check if Docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("Docker is not installed: %w", err)
	}
	color.Green("Docker is installed.")

	return nil
}

func getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile string, step int) error {
	color.Yellow("Extracting passwords for stats using --show...")

	currentCount, err := utils.CountLines(cumulativeCrackedFile)
	if err != nil {
		return err
	}

	var hashcatCommand []string
	hashcatCommand = []string{"-m", hashcatMode, hashlist, "--show"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommandToFile(hashcatPath, hashcatCommand, cumulativeCrackedFile) // Ignore error as Hashcat may return 1 even on success

	newCount, err := utils.CountLines(cumulativeCrackedFile)
	if err != nil {
		return err
	}

	color.Red("Extracted %d new passwords for stats.", newCount-currentCount)

	statsMessage := fmt.Sprintf("Extracted %d new passwords for step %d.", newCount-currentCount, step)

	if err := utils.AppendToFile(cumulativeCrackedStatsFile, []string{statsMessage}); err != nil {
		return fmt.Errorf("error writing cracked passwords to file: %w", err)
	}

	return nil
}

func ProcessHashcatTasks(hashlist, wordlist, potfile, clemRule, rulesFull, cewlURL, cewlWordlist, hashcatPath, hashcatMode, passphrases, passphraseRule1, passphraseRule2, dictionary string, enableAdditionalWordlists bool) error {

	// Step 1: Validate hashlist
	color.Yellow("Validating hashlist...")
	if err := utils.ValidateFileExists(hashlist); err != nil {
		return fmt.Errorf("hashlist validation failed: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	cumulativeCrackedFile := filepath.Join(config.CacheDir, fmt.Sprintf("cumulative_cracked_%s.txt", timestamp))
	cumulativeCrackedStatsFile := filepath.Join(config.CacheDir, fmt.Sprintf("cumulative_cracked_stats_%s.txt", timestamp))

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 1)

	// Step 2: Extract passwords using --show and process them
	color.Yellow("Extracting passwords using --show...")
	tempCrackedFile := filepath.Join(config.CacheDir, fmt.Sprintf("temp_cracked_passwords_%s.txt", timestamp))
	var hashcatCommand []string

	hashcatCommand = []string{"-m", hashcatMode, hashlist, "--show"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommandToFile(hashcatPath, hashcatCommand, tempCrackedFile) // Ignore error as Hashcat may return 1 even on success
	passwords, err := utils.ExtractPasswords(tempCrackedFile)
	if err != nil {
		return fmt.Errorf("error processing cracked passwords: %w", err)
	}
	passwordsFile := filepath.Join(config.CacheDir, fmt.Sprintf("cracked_passwords_%s.txt", timestamp))
	if err := utils.WriteToFile(passwordsFile, passwords); err != nil {
		return fmt.Errorf("error writing cracked passwords to file: %w", err)
	}
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, passwordsFile, "-r", rulesFull, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Cracked password processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 2)

	// Step 3: Get passwords with custom potfile and process them
	color.Yellow("Extracting passwords from custom potfile using --show...")
	tempCrackedFile = filepath.Join(config.CacheDir, fmt.Sprintf("temp_custom_potfile_cracked_passwords_%s.txt", timestamp))

	hashcatCommand = []string{"-m", hashcatMode, hashlist, "--show", "--potfile-path", potfile}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommandToFile(hashcatPath, hashcatCommand, tempCrackedFile) // Ignore error as Hashcat may return 1 even on success
	passwords, err = utils.ExtractPasswords(tempCrackedFile)
	if err != nil {
		return fmt.Errorf("error processing cracked passwords: %w", err)
	}
	passwordsFile = filepath.Join(config.CacheDir, fmt.Sprintf("custom_potfile_cracked_passwords_%s.txt", timestamp))
	if err := utils.WriteToFile(passwordsFile, passwords); err != nil {
		return fmt.Errorf("error writing cracked passwords to file: %w", err)
	}
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, passwordsFile, "-r", rulesFull, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Cracked password from potfile processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 3)

	// Step 4: Run rockyou.txt wordlist
	color.Yellow("Running rockyou.txt wordlist...")
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, wordlist, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Wordlist processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 4)

	// Step 5: Run rockyou.txt with clem9669_large.rule
	color.Yellow("Running rockyou.txt with clem9669_large.rule...")
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, wordlist, "-r", clemRule, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Rule-based processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 5)

	// Step 6: Extract usernames and run with rules_full.rule
	color.Yellow("Extracting usernames and running with rules_full.rule...")
	usernames, err := utils.ExtractUsernames(hashlist)
	if err != nil {
		return fmt.Errorf("error extracting usernames: %w", err)
	}

	usernameFile := filepath.Join(config.CacheDir, fmt.Sprintf("usernames_%s.txt", timestamp))
	if err := utils.WriteToFile(usernameFile, usernames); err != nil {
		return fmt.Errorf("error writing usernames to file: %w", err)
	}
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, usernameFile, "-r", rulesFull, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Username-based processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 6)

	// Step 7: Use CeWL to generate a wordlist and run with rules_full.rule
	if cewlURL != "" {
		cewlOutputFile := filepath.Join(config.CacheDir, fmt.Sprintf("cewl_wordlist_%s.txt", timestamp))
		dockerCommand := []string{
			"run", "--rm",
			"-v", fmt.Sprintf("%s:/output", config.CacheDir),
			"ghcr.io/digininja/cewl",
			"-w", fmt.Sprintf("/output/cewl_wordlist_%s.txt", timestamp), cewlURL,
			"--with-numbers", "--meta", "--email"}
		fmt.Printf("DEBUG: Docker command: docker %v\n", dockerCommand)
		_ = utils.RunCommand("docker", dockerCommand) // Ignore error as Docker may return non-critical errors

		// Clean the generated CeWL wordlist
		color.Yellow("Cleaning CeWL wordlist to remove non-ASCII characters...")
		cleanedWordlist, err := utils.CleanWordlist(cewlOutputFile, config.CacheDir, timestamp)
		if err != nil {
			return fmt.Errorf("failed to clean CeWL wordlist: %w", err)
		}
		color.Green("CeWL wordlist cleaned successfully: %s", cleanedWordlist)

		// Use the cleaned wordlist with Hashcat
		hashcatCommand := []string{"-a", "0", "-m", hashcatMode, hashlist, cleanedWordlist, "-r", rulesFull, "--status", "--status-timer", "30"}
		fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
		_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
		color.Green("CeWL-based processing completed.")

		getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 7)
	}

	// Step 8: Process passphrases with two rules
	color.Yellow("Processing passphrases with two rules...")
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, passphrases, "-r", passphraseRule1, "-r", passphraseRule2, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Passphrase processing completed.")

	// Step 9: Extract and process additional wordlists
	if enableAdditionalWordlists {
		color.Yellow("Processing additional wordlists with Hashcat...")
		extraWordlistCommands := []string{
			"7z x -so /media/extra2/Wordlists/all_in_one.txt.7z | /media/extra/Work/tools/hashcat-6.2.6/hashcat.bin -a 0 -m " + hashcatMode + " " + hashlist,
			"7z x -so /media/extra2/Wordlists/hashmob.net_2024-12-01.found.7z | /media/extra/Work/tools/hashcat-6.2.6/hashcat.bin -a 0 -m " + hashcatMode + " " + hashlist,
			"bzip2 -dc /media/extra2/Wordlists/rockyou2024.txt.bz2 | /media/extra/Work/tools/hashcat-6.2.6/hashcat.bin -a 0 -m " + hashcatMode + " " + hashlist,
		}

		for _, cmd := range extraWordlistCommands {
			fmt.Printf("DEBUG: Running shell command: %s\n", cmd)
			_ = utils.RunShellCommand(cmd) // Use a utility to run shell commands
		}
		color.Green("Additional wordlists processed.")

		getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 8)
	}

	// Step 9: Process dictionary with rules_full.rule
	color.Yellow("Running dictonary with rules_full.rule...")
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, dictionary, "-r", rulesFull, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Dictionary processing completed.")

	// Step 10: Extract passwords using --show and process them
	color.Yellow("Extracting passwords using --show...")
	tempCrackedFile = filepath.Join(config.CacheDir, fmt.Sprintf("temp_cracked_passwords_%s.txt", timestamp))

	hashcatCommand = []string{"-m", hashcatMode, hashlist, "--show"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommandToFile(hashcatPath, hashcatCommand, tempCrackedFile) // Ignore error as Hashcat may return 1 even on success
	passwords, err = utils.ExtractPasswords(tempCrackedFile)
	if err != nil {
		return fmt.Errorf("error processing cracked passwords: %w", err)
	}
	passwordsFile = filepath.Join(config.CacheDir, fmt.Sprintf("cracked_passwords_%s.txt", timestamp))
	if err := utils.WriteToFile(passwordsFile, passwords); err != nil {
		return fmt.Errorf("error writing cracked passwords to file: %w", err)
	}
	hashcatCommand = []string{"-a", "0", "-m", hashcatMode, hashlist, passwordsFile, "-r", rulesFull, "--status", "--status-timer", "30"}
	fmt.Printf("DEBUG: Hashcat command: %s %v\n", hashcatPath, hashcatCommand)
	_ = utils.RunCommand(hashcatPath, hashcatCommand) // Ignore error as Hashcat may return 1 even on success
	color.Green("Cracked password processing completed.")

	getPasswordStats(hashcatMode, hashcatPath, hashlist, cumulativeCrackedFile, cumulativeCrackedStatsFile, 9)

	color.Green("All steps completed successfully.")
	return nil
}
