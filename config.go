package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = ".gogrc"

// getConfigFilePath returns the path to the config file where the GitHub username and shortcuts are stored.
func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, configFileName)
}

// readConfig reads the GitHub username and shortcuts from the config file.
func readConfig() (string, map[string]string, error) {
	configFilePath := getConfigFilePath()
	file, err := os.Open(configFilePath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	var username string
	shortcuts := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		} else if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				shortcuts[parts[0]] = parts[1]
			}
		}
	}
	return username, shortcuts, nil
}

// writeConfig writes the GitHub username and shortcuts to the config file.
func writeConfig(username string, shortcuts map[string]string) error {
	configFilePath := getConfigFilePath()
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if username != "" {
		_, _ = file.WriteString("username=" + username + "\n")
	}
	for key, value := range shortcuts {
		_, _ = file.WriteString(key + "=" + value + "\n")
	}
	return nil
}
