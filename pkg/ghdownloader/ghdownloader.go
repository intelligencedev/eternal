package ghdownloader

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var HttpGet = http.Get

type RepoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func DownloadRepo(owner, repo, branch, destPath string) error {
	if branch == "" {
		var err error
		branch, err = getDefaultBranch(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get default branch: %w", err)
		}
	}

	url := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", owner, repo, branch)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download repository: HTTP status %d", resp.StatusCode)
	}

	zipPath := filepath.Join(destPath, fmt.Sprintf("%s-%s.zip", repo, branch))
	out, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write zip file: %w", err)
	}

	return nil
}

func getDefaultBranch(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get repository info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get repository info: HTTP status %d", resp.StatusCode)
	}

	var repoInfo RepoInfo
	err = json.NewDecoder(resp.Body).Decode(&repoInfo)
	if err != nil {
		return "", fmt.Errorf("failed to decode repository info: %w", err)
	}

	return repoInfo.DefaultBranch, nil
}

func ExtractZip(zipPath, destPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

func DownloadAndExtractRepo(owner, repo, branch, destPath string) error {
	// Create a temporary directory for the zip file
	tempDir, err := os.MkdirTemp("", "repo-download-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download the repository
	err = DownloadRepo(owner, repo, branch, tempDir)
	if err != nil {
		return err
	}

	// Get the path of the downloaded zip file
	zipFiles, err := filepath.Glob(filepath.Join(tempDir, "*.zip"))
	if err != nil || len(zipFiles) == 0 {
		return fmt.Errorf("failed to find downloaded zip file: %w", err)
	}
	zipPath := zipFiles[0]

	// Extract the zip file
	err = ExtractZip(zipPath, destPath)
	if err != nil {
		return err
	}

	return nil
}
