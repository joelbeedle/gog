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

// getConfigFilePath returns the path to the config file where the GitHub username is stored.
func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, configFileName)
}

// readUsernameFromFile reads the GitHub username from the config file.
func readUsernameFromFile() (string, error) {
	configFilePath := getConfigFilePath()
	file, err := os.Open(configFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("no username found")
}

// writeUsernameToFile writes the GitHub username to the config file.
func writeUsernameToFile(username string) error {
	configFilePath := getConfigFilePath()
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(username)
	return err
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
	var err error

	// Check if the user wants to update the username
	if len(os.Args) > 1 && os.Args[1] == "-u" && len(os.Args) > 2 {
		newUsername := os.Args[2]
		err = writeUsernameToFile(newUsername)
		if err != nil {
			fmt.Println("Error saving username:", err)
			return
		}
		fmt.Println("GitHub username updated successfully.")
		return
	}

	// Try to read the username from the config file
	username, err = readUsernameFromFile()
	if err != nil {
		fmt.Println("No GitHub username found.")
		// Prompt the user to enter a GitHub username
		username = promptUsername()
		err = writeUsernameToFile(username)
		if err != nil {
			fmt.Println("Error saving username:", err)
			return
		}
		fmt.Println("GitHub username saved successfully.")
	}

	// Use the username to construct the URL
	var repo string
	if len(os.Args) > 1 {
		repo = os.Args[1]
	} else {
		repo = ""
	}
	url := constructURL(username, repo)

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
