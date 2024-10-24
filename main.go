package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	var repo string

	// Check if an argument is provided
	if len(os.Args) > 1 {
		repo = os.Args[1] // Use the provided argument
	} else {
		repo = "" // Default to your profile if no argument is provided
	}

	// Construct the URL
	url := constructURL(repo)

	// Open the URL in the default browser
	err := openBrowser(url)
	if err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}

// constructURL builds the URL based on the repository or defaults to your profile
func constructURL(repo string) string {
	base := "https://github.com/joelbeedle"
	if repo == "" {
		return base // Open profile if no repo is provided
	}
	return fmt.Sprintf("%s/%s", base, repo) // Or open repo if provided
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
