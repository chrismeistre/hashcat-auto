package main

import (
	"flag"
	"fmt"
	"hashcat-auto/cmd"
	"hashcat-auto/config"
	"hashcat-auto/utils"
	"os"
	"os/exec"
	"path/filepath"

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

func main() {
	// Define the path to config.json
	configPath := filepath.Join("config.json")

	// Load configuration
	err := config.LoadConfig(configPath)
	if err != nil {
		color.Red("Error loading config: %v", err)
		os.Exit(1)
	}

	// Print the loaded configuration for debugging
	color.Green("Loaded Configuration:")
	color.Green("DefaultHashcatPath: %s\n", config.DefaultHashcatPath)
	color.Green("DefaultWordlist: %s\n", config.DefaultWordlist)
	color.Green("DefaultPotfile: %s\n", config.DefaultPotfile)
	color.Green("DefaultClemRule: %s\n", config.DefaultClemRule)
	color.Green("DefaultRulesFull: %s\n", config.DefaultRulesFull)
	color.Green("DefaultPassphrases: %s\n", config.DefaultPassphrases)
	color.Green("DefaultPassphraseRule1: %s\n", config.DefaultPassphraseRule1)
	color.Green("DefaultPassphraseRule2: %s\n", config.DefaultPassphraseRule2)
	color.Green("DefaultDictionary: %s\n", config.DefaultDictionary)
	color.Green("Cache Directory: %s\n", config.CacheDir)

	// Define command-line flags
	hashlist := flag.String("hashlist", "", "Path to the hashlist file (REQUIRED: user:hash format)")
	wordlist := flag.String("wordlist", config.DefaultWordlist, "Path to the wordlist file")
	potfile := flag.String("potfile", config.DefaultPotfile, "Path to the potfile file")
	clemRule := flag.String("clemrule", config.DefaultClemRule, "Path to clem9669_large.rule file")
	rulesFull := flag.String("rulesfull", config.DefaultRulesFull, "Path to rules_full.rule file")
	cewlURL := flag.String("url", "", "URL for CeWL to generate wordlist")
	cewlWordlist := flag.String("cewlwordlist", "cewl_wordlist.txt", "Output file for CeWL wordlist")
	hashcatPath := flag.String("hashcat", config.DefaultHashcatPath, "Path to the hashcat binary")
	hashcatMode := flag.String("mode", "", "Hashcat mode to use (REQUIRED)")
	passphrases := flag.String("passphrases", config.DefaultPassphrases, "Path to passphrases wordlist")
	passphraseRule1 := flag.String("passphraserule1", config.DefaultPassphraseRule1, "Path to passphrase-rule1.rule")
	passphraseRule2 := flag.String("passphraserule2", config.DefaultPassphraseRule2, "Path to passphrase-rule2.rule")
	enableAdditionalWordlists := flag.Bool("enable-additional-wordlists", false, "Enable processing of additional wordlists")
	dictionary := flag.String("dictionary", config.DefaultDictionary, "Path to the dictionary file")

	// Parse command-line flags
	flag.Parse()

	// Validate required flags
	if *hashlist == "" {
		color.Red("Error: --hashlist is required")
		flag.Usage()
		os.Exit(1)
	}

	if *hashcatMode == "" {
		color.Red("Error: --mode is required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate environment and input files
	err = validateEnvironment(*hashcatPath)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	filesToValidate := []string{
		*hashlist,
		*wordlist,
		*potfile,
		*clemRule,
		*rulesFull,
		*passphrases,
		*passphraseRule1,
	}

	if err := utils.ValidateFiles(filesToValidate); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	// Run Hashcat tasks
	if err := cmd.ProcessHashcatTasks(*hashlist, *wordlist, *potfile, *clemRule, *rulesFull, *cewlURL, *cewlWordlist, *hashcatPath, *hashcatMode, *passphrases, *passphraseRule1, *passphraseRule2, *dictionary, *enableAdditionalWordlists); err != nil {
		color.Red("Error: %v", err)
		return
	}

	color.Green("All tasks completed successfully.")
}
