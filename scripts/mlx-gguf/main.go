package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// Define the parameters for the Python command
	prompt := "Write a comprehensive plan to become an expert software systems architect in markdown format."
	pythonCommand := "python3"
	scriptPath := "generate.py"
	modelPath := "/Users/art/.eternal/data/models/Zephyr-7B-Beta/zephyr-7b-beta.Q8_0.gguf"
	promptTemplate := fmt.Sprintf("<|system|></s><|user|>%s</s><|assistant|>", prompt)
	maxTokens := "8092"

	// Run the Python command and stream the response
	err := StreamPythonCommand(pythonCommand, scriptPath, modelPath, promptTemplate, maxTokens)
	if err != nil {
		log.Fatalf("Failed to run python command: %v", err)
	}
}

func StreamPythonCommand(pythonCommand, scriptPath, modelPath, promptTemplate, maxTokens string) error {
	// Construct the command with arguments
	cmd := exec.Command(pythonCommand, scriptPath, "--gguf", modelPath, "--prompt", promptTemplate, "--max-tokens", maxTokens)

	// Create a pipe for the output of the script
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a channel to receive the output
	outputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
		close(outputChan)
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Println("Error:", scanner.Text())
		}
	}()

	for line := range outputChan {
		fmt.Println(line)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
