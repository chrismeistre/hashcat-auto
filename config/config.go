package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Declare global variables
var (
	DefaultHashcatPath     string
	DefaultWordlist        string
	DefaultPotfile         string
	DefaultClemRule        string
	DefaultRulesFull       string
	DefaultPassphrases     string
	DefaultPassphraseRule1 string
	DefaultPassphraseRule2 string
	DefaultDictionary      string
	CacheDir               string
)

// Config struct to map JSON keys
type Config struct {
	HashcatPath     string `json:"hashcat_path"`
	Wordlist        string `json:"wordlist"`
	Potfile         string `json:"potfile"`
	ClemRule        string `json:"clem_rule"`
	RulesFull       string `json:"rules_full"`
	Passphrases     string `json:"passphrases"`
	PassphraseRule1 string `json:"passphrase_rule1"`
	PassphraseRule2 string `json:"passphrase_rule2"`
	Dictionary      string `json:"dictionary"`
	CacheDir        string `json:"cache_dir"`
}

// LoadConfig reads the config.json file and assigns values to global variables
func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	// Assign values to global variables
	DefaultHashcatPath = cfg.HashcatPath
	DefaultWordlist = cfg.Wordlist
	DefaultPotfile = cfg.Potfile
	DefaultClemRule = cfg.ClemRule
	DefaultRulesFull = cfg.RulesFull
	DefaultPassphrases = cfg.Passphrases
	DefaultPassphraseRule1 = cfg.PassphraseRule1
	DefaultPassphraseRule2 = cfg.PassphraseRule2
	DefaultDictionary = cfg.Dictionary
	CacheDir = cfg.CacheDir

	// Ensure cache directory exists
	if err := os.MkdirAll(CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	return nil
}
