package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Prompts the user to enter a GitHub username.
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

	getConfigCmd := &cobra.Command{
		Use:   "get-config",
		Short: "Gets the contents of the config file",
		Run: func(cmd *cobra.Command, args []string) {
			configContents, err := getConfig()
			if err != nil {
				fmt.Println("Error getting config:", err)
				return
			}

			for i, line := range configContents {
				fmt.Printf("%d: %s\n", i+1, line)
			}
		},
	}

	rootCmd.AddCommand(setUsernameCmd, addShortcutCmd, getConfigCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
