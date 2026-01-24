package ollama

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// InstallResult contains the result of an installation attempt.
type InstallResult struct {
	Success bool
	Message string
	Output  string
}

// Install attempts to install Ollama on the current system.
// Returns an InstallResult with status and any output/errors.
func Install(ctx context.Context) InstallResult {
	switch runtime.GOOS {
	case "darwin":
		return installMacOS(ctx)
	case "linux":
		return installLinux(ctx)
	default:
		return InstallResult{
			Success: false,
			Message: fmt.Sprintf("Auto-install is not supported on %s. Please install manually from https://ollama.ai/download", runtime.GOOS),
		}
	}
}

// installMacOS installs Ollama using Homebrew on macOS.
func installMacOS(ctx context.Context) InstallResult {
	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return InstallResult{
			Success: false,
			Message: "Homebrew is not installed. Please install Homebrew first or download Ollama from https://ollama.ai/download/mac",
		}
	}

	// Install using Homebrew
	cmd := exec.CommandContext(ctx, "brew", "install", "ollama")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Check if it's already installed
		if strings.Contains(outputStr, "already installed") {
			return InstallResult{
				Success: true,
				Message: "Ollama is already installed",
				Output:  outputStr,
			}
		}
		return InstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to install Ollama: %v", err),
			Output:  outputStr,
		}
	}

	// Verify installation
	if IsInstalled() {
		return InstallResult{
			Success: true,
			Message: "Ollama installed successfully",
			Output:  outputStr,
		}
	}

	return InstallResult{
		Success: false,
		Message: "Installation completed but Ollama binary not found",
		Output:  outputStr,
	}
}

// installLinux installs Ollama using the official install script on Linux.
func installLinux(ctx context.Context) InstallResult {
	// Check if curl is available
	if _, err := exec.LookPath("curl"); err != nil {
		return InstallResult{
			Success: false,
			Message: "curl is not installed. Please install curl first or download Ollama from https://ollama.ai/download/linux",
		}
	}

	// Use the official install script
	cmd := exec.CommandContext(ctx, "sh", "-c", "curl -fsSL https://ollama.ai/install.sh | sh")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return InstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to install Ollama: %v", err),
			Output:  outputStr,
		}
	}

	// Verify installation
	if IsInstalled() {
		return InstallResult{
			Success: true,
			Message: "Ollama installed successfully",
			Output:  outputStr,
		}
	}

	return InstallResult{
		Success: false,
		Message: "Installation completed but Ollama binary not found",
		Output:  outputStr,
	}
}

// CanAutoInstall returns whether auto-installation is supported on the current platform.
func CanAutoInstall() bool {
	switch runtime.GOOS {
	case "darwin":
		// Check if Homebrew is available
		_, err := exec.LookPath("brew")
		return err == nil
	case "linux":
		// Check if curl is available
		_, err := exec.LookPath("curl")
		return err == nil
	default:
		return false
	}
}

// InstallMethod returns a description of how Ollama will be installed.
func InstallMethod() string {
	switch runtime.GOOS {
	case "darwin":
		return "Install via Homebrew (brew install ollama)"
	case "linux":
		return "Install via official script (curl -fsSL https://ollama.ai/install.sh | sh)"
	default:
		return "Manual installation required"
	}
}
