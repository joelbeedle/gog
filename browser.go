package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Builds the URL based on the repository or defaults to the user's profile
func constructURL(username, repo string) string {
	base := fmt.Sprintf("https://github.com/%s", username)
	if repo == "" {
		return base // Open profile if no repo is provided
	}
	return fmt.Sprintf("%s/%s", base, repo) // Open repo if provided
}

// Opens the specified URL in the default browser
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
