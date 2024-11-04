package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const configFileName = ".gog_config"

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

// promptUsername prompts the user to enter a GitHub username.
func promptUsername() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your GitHub username: ")
	username, _ := reader.ReadString('\n')
	return strings.TrimSpace(username)
}

func main() {
	var username string
	var shortcuts map[string]string
	var err error

	// Check if the user wants to update the username or set a shortcut
	if len(os.Args) > 1 && os.Args[1] == "-u" && len(os.Args) > 2 {
		newUsername := os.Args[2]
		username, shortcuts, err = readConfig()
		if err != nil {
			username, shortcuts = "", make(map[string]string)
		}
		username = newUsername
		err = writeConfig(username, shortcuts)
		if err != nil {
			fmt.Println("Error saving username:", err)
			return
		}
		fmt.Println("GitHub username updated successfully.")
		return
	} else if len(os.Args) > 1 && os.Args[1] == "-s" && len(os.Args) > 3 {
		key := os.Args[2]
		value := os.Args[3]
		username, shortcuts, err = readConfig()
		if err != nil {
			username, shortcuts = "", make(map[string]string)
		}
		shortcuts[key] = value
		err = writeConfig(username, shortcuts)
		if err != nil {
			fmt.Println("Error saving shortcut:", err)
			return
		}
		fmt.Println("Shortcut added successfully.")
		return
	}

	// Try to read the username and shortcuts from the config file
	username, shortcuts, err = readConfig()
	if err != nil || username == "" {
		fmt.Println("No GitHub username found.")
		// Prompt the user to enter a GitHub username
		username = promptUsername()
		err = writeConfig(username, shortcuts)
		if err != nil {
			fmt.Println("Error saving username:", err)
			return
		}
		fmt.Println("GitHub username saved successfully.")
	}

	// Handle repository or shortcut arguments
	var url string
	if len(os.Args) > 1 {
		if len(os.Args) == 2 {
			// Check if it's a shortcut
			if val, exists := shortcuts[os.Args[1]]; exists {
				url = fmt.Sprintf("https://github.com/%s", val)
			} else {
				url = constructURL(username, os.Args[1])
			}
		} else {
			// Multiple arguments for custom user/repo
			url = fmt.Sprintf("https://github.com/%s/%s", os.Args[1], os.Args[2])
		}
	} else {
		// Default to the user's profile
		url = fmt.Sprintf("https://github.com/%s", username)
	}

	// Open the URL in the default browser
	err = openBrowser(url)
	if err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}

// constructURL builds the URL based on the repository or defaults to the user's profile
func constructURL(username, repo string) string {
	base := fmt.Sprintf("https://github.com/%s", username)
	if repo == "" {
		return base // Open profile if no repo is provided
	}
	return fmt.Sprintf("%s/%s", base, repo) // Open repo if provided
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	// Handle cross-platform behavior
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
