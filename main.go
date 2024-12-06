package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// promptUsername prompts the user to enter a GitHub username.
func promptUsername() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your GitHub username: ")
	username, _ := reader.ReadString('\n')
	return strings.TrimSpace(username)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gog [shortcut or repo]",
		Short: "GitHub CLI tool for quickly opening repositories and profiles",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Default behavior: Handle repository or shortcut arguments
			username, shortcuts, _ := readConfig()
			var url string
			if len(args) == 1 {
				if val, exists := shortcuts[args[0]]; exists {
					url = fmt.Sprintf("https://github.com/%s", val)
				} else {
					url = constructURL(username, args[0])
				}
			} else if len(args) == 2 {
				url = fmt.Sprintf("https://github.com/%s/%s", args[0], args[1])
			} else {
				url = fmt.Sprintf("https://github.com/%s", username)
			}
			err := openBrowser(url)
			if err != nil {
				fmt.Println("Failed to open browser:", err)
			}
		},
	}

	// Add `set-username` command
	setUsernameCmd := &cobra.Command{
		Use:   "set-username [username]",
		Short: "Set your GitHub username",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			newUsername := args[0]
			_, shortcuts, _ := readConfig()
			err := writeConfig(newUsername, shortcuts)
			if err != nil {
				fmt.Println("Error saving username:", err)
				return
			}
			fmt.Println("GitHub username updated successfully.")
		},
	}

	// Add `add-shortcut` command
	addShortcutCmd := &cobra.Command{
		Use:   "add-shortcut [key] [repo]",
		Short: "Add a shortcut for a repository",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			username, shortcuts, _ := readConfig()
			shortcuts[key] = value
			err := writeConfig(username, shortcuts)
			if err != nil {
				fmt.Println("Error saving shortcut:", err)
				return
			}
			fmt.Println("Shortcut added successfully.")
		},
	}

	rootCmd.AddCommand(setUsernameCmd, addShortcutCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
