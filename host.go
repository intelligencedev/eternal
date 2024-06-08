// NOTE: Thiese functions are not implemented in the main app yet.
package main

import (
	// Importing necessary packages for executing commands, formatting strings, and hardware information retrieval.
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	// Package for hardware information
	"github.com/shirou/gopsutil/mem" // Package for system memory information
)

// HostInfo struct: Stores information about the host system.
type HostInfo struct {
	OS     string `json:"os"`   // Operating System
	Arch   string `json:"arch"` // Architecture (e.g., amd64, 386)
	CPUs   int    `json:"cpus"` // Number of CPUs
	Memory struct {
		Total uint64 `json:"total"` // Total memory in bytes
	} `json:"memory"`
	GPUs []GPUInfo `json:"gpus"` // Slice of GPU information
}

// GPUInfo struct: Stores information about GPUs in the system.
type GPUInfo struct {
	Model              string `json:"model"`                 // GPU model
	TotalNumberOfCores string `json:"total_number_of_cores"` // Total cores in GPU
	MetalSupport       string `json:"metal_support"`         // Metal support (specific to macOS)
}

// GetHostInfo function: Retrieves information about the host system.
func GetHostInfo() (HostInfo, error) {
	hostInfo := HostInfo{
		OS:   runtime.GOOS,     // Fetching OS
		Arch: runtime.GOARCH,   // Fetching architecture
		CPUs: runtime.NumCPU(), // Fetching CPU count
	}

	// Retrieve memory information using gopsutil
	vmStat, _ := mem.VirtualMemory()
	hostInfo.Memory.Total = vmStat.Total

	// GPU information retrieval based on OS
	switch runtime.GOOS {
	case "darwin":
		// macOS specific GPU information retrieval
		gpus, err := getMacOSGPUInfo()
		if err != nil {
			fmt.Printf("Error getting GPU info: %v\n", err)
		} else {
			hostInfo.GPUs = append(hostInfo.GPUs, gpus)
		}

	case "linux", "windows":
		// Disabling since this needs more work
		// Linux and Windows GPU information retrieval
		// gpu, err := ghw.GPU()
		// if err != nil {
		// 	fmt.Printf("Error getting GPU info: %v\n", err)
		// } else {
		// 	for _, card := range gpu.GraphicsCards {
		// 		gpuInfo := GPUInfo{
		// 			Model: card.DeviceInfo.Product.Name, // Fetching GPU model
		// 		}
		// 		hostInfo.GPUs = append(hostInfo.GPUs, gpuInfo)
		// 	}
		// }
	}

	return hostInfo, nil
}

// getMacOSGPUInfo function: Retrieves GPU information for macOS.
func getMacOSGPUInfo() (GPUInfo, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return GPUInfo{}, err
	}

	return parseGPUInfo(out.String())
}

// parseGPUInfo function: Parses the output from system_profiler to extract GPU info.
func parseGPUInfo(input string) (GPUInfo, error) {
	gpuInfo := GPUInfo{}

	for _, line := range strings.Split(input, "\n") {
		// Extracting relevant information from the output
		if strings.Contains(line, "Chipset Model") {
			gpuInfo.Model = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Total Number of Cores") {
			gpuInfo.TotalNumberOfCores = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Metal") {
			gpuInfo.MetalSupport = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return gpuInfo, nil
}
