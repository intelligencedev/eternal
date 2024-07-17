// comfyui.go
package main

import (
	"bufio"
	"context"
	"errors"
	"eternal/internal/python"
	"eternal/pkg/ghdownloader"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pterm/pterm"
)

// startComfyUI initializes and starts the ComfyUI service.
// It installs necessary requirements, runs the service, and waits for it to be ready.
func startComfyUI(ctx context.Context, config *AppConfig) error {
	comfyPort := config.ServiceHosts["image"]["image_host_1"].Port
	comfyUIPath := filepath.Join(config.DataPath, "sd/ComfyUI-master")

	if err := installComfyUIRequirements(comfyUIPath); err != nil {
		return err
	}

	if err := runComfyUI(ctx, comfyUIPath, comfyPort); err != nil {
		return err
	}

	return waitForComfyUI(config)
}

// installComfyUIRequirements installs the necessary Python packages for ComfyUI.
func installComfyUIRequirements(comfyUIPath string) error {
	requirementsPath := filepath.Join(comfyUIPath, "requirements.txt")
	output, err := python.ExecuteScript("-m", "pip", "install", "-r", requirementsPath)
	if err != nil {
		return fmt.Errorf("failed to install ComfyUI requirements: %v\nOutput: %s", err, output)
	}
	pterm.Info.Println("ComfyUI requirements installed")
	return nil
}

// runComfyUI starts the ComfyUI service with the specified configuration.
func runComfyUI(ctx context.Context, comfyUIPath, comfyPort string) error {
	cmdArgs := []string{
		filepath.Join(comfyUIPath, "main.py"),
		"--listen",
		"--port", comfyPort,
		"--disable-auto-launch",
		"--preview-method", "none",
		//"--highvram",
		"--dont-print-server",
		//"--force-fp16",
		//"--use-split-cross-attention",
	}

	comfyUICmd = exec.CommandContext(ctx, "python", cmdArgs...)
	comfyUICmd.Dir = comfyUIPath

	stdout, err := comfyUICmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := comfyUICmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	if err := comfyUICmd.Start(); err != nil {
		return fmt.Errorf("failed to start ComfyUI: %v", err)
	}

	go monitorComfyUIOutput(stdout, stderr)

	return nil
}

// monitorComfyUIOutput monitors and logs the output of the ComfyUI service.
func monitorComfyUIOutput(stdout, stderr io.Reader) {
	go scanAndLog(stdout, "ComfyUI: ", log.Info)
	go scanAndLog(stderr, "ComfyUI: ", log.Info)
}

// scanAndLog reads and logs the output from the provided reader.
func scanAndLog(r io.Reader, prefix string, logFunc func(...interface{})) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logFunc(prefix, scanner.Text())
	}
}

// waitForComfyUI waits for the ComfyUI service to become available.
func waitForComfyUI(config *AppConfig) error {
	comfyHost := config.ServiceHosts["image"]["image_host_1"].Host
	comfyPort := config.ServiceHosts["image"]["image_host_1"].Port
	comfyUrl := fmt.Sprintf("http://%s:%s", comfyHost, comfyPort)
	client := &http.Client{Timeout: 1 * time.Second}
	for i := 0; i < 120; i++ {
		if resp, err := client.Get(comfyUrl); err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for ComfyUI to start")
}

// stopComfyUI stops the ComfyUI service gracefully.
func stopComfyUI(ctx context.Context) error {
	if comfyUICmd == nil || comfyUICmd.Process == nil {
		return nil
	}

	// Send SIGTERM to ComfyUI
	if err := comfyUICmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to ComfyUI: %v", err)
	}

	// Create a channel to signal when ComfyUI has exited
	done := make(chan error, 1)
	go func() {
		done <- comfyUICmd.Wait()
	}()

	// Wait for ComfyUI to exit or for the context to be canceled
	select {
	case <-ctx.Done():
		// Context was canceled (e.g., due to a timeout or signal)
		if err := comfyUICmd.Process.Kill(); err != nil {
			return nil
		}
		return nil
	case <-time.After(10 * time.Second):
		// Timeout reached
		if err := comfyUICmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill ComfyUI process: %v", err)
		}
		return nil
	case err := <-done:
		// ComfyUI exited
		if err != nil {
			return fmt.Errorf("ComfyUI exited with error: %v", err)
		}
		return nil
	}
}

// setupComfyUI sets up the ComfyUI directory and installs requirements.
func setupComfyUI(configPath string) error {
	comfyuiPath := filepath.Join(configPath, "sd/ComfyUI-master")
	if _, err := os.Stat(comfyuiPath); os.IsNotExist(err) {
		if err := ghdownloader.DownloadAndExtractRepo("intelligencedev", "ComfyUI", "", filepath.Join(configPath, "sd")); err != nil {
			return fmt.Errorf("failed to download ComfyUI: %w", err)
		}

		if err := installPythonRequirements(filepath.Join(comfyuiPath, "requirements.txt")); err != nil {
			return fmt.Errorf("failed to install ComfyUI requirements: %w", err)
		}
	}
	return nil
}

// setupImpactPack sets up the Impact Pack for ComfyUI.
func setupImpactPack(configPath string) error {
	impactPackPath := filepath.Join(configPath, "sd/ComfyUI-master/custom_nodes/ComfyUI-Impact-Pack-Main")
	if _, err := os.Stat(impactPackPath); os.IsNotExist(err) {
		if err := ghdownloader.DownloadAndExtractRepo("ltdrdata", "ComfyUI-Impact-Pack", "", filepath.Join(configPath, "sd/ComfyUI-master/custom_nodes")); err != nil {
			return fmt.Errorf("failed to download Impact Pack: %w", err)
		}

		if err := installPythonRequirements(filepath.Join(impactPackPath, "requirements.txt")); err != nil {
			return fmt.Errorf("failed to install Impact Pack requirements: %w", err)
		}
	}
	return nil
}

// setupKolors sets up the Kolors library for ComfyUI.
func setupKolors(configPath string) error {
	kolorsPath := filepath.Join(configPath, "sd/ComfyUI-master/custom_nodes/ComfyUI-KwaiKolorsWrapper-main")
	if _, err := os.Stat(kolorsPath); os.IsNotExist(err) {
		if err := ghdownloader.DownloadAndExtractRepo("kijai", "ComfyUI-KwaiKolorsWrapper", "", filepath.Join(configPath, "sd/ComfyUI-master/custom_nodes")); err != nil {
			return fmt.Errorf("failed to download Kolors: %w", err)
		}

		if err := installPythonRequirements(filepath.Join(kolorsPath, "requirements.txt")); err != nil {
			return fmt.Errorf("failed to install Kolors requirements: %w", err)
		}
	}

	output, err := python.ExecuteScript("-m", "pip", "install", "sentencepiece")
	if err != nil {
		return fmt.Errorf("failed to install ComfyUI requirements: %v\nOutput: %s", err, output)
	}

	return nil
}
