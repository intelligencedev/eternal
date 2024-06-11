package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// Define the parameters for the CLI command
	cliCommand := "mlx_lm.generate"
	//model := "mlx-community/Phi-3-mini-128k-instruct-8bit"
	model := "/Users/arturoaquino/.eternal-v1/models/Phi-3-medium-128k-instruct"
	prompt := "write a quicksort in python."
	maxTokens := "32000"

	// Run the CLI command and stream the response
	err := StreamCLICommand(cliCommand, model, prompt, maxTokens)
	if err != nil {
		log.Fatalf("Failed to run CLI command: %v", err)
	}
}

func StreamCLICommand(cliCommand, model, prompt, maxTokens string) error {
	// Construct the command with arguments
	cmd := exec.Command(cliCommand, "--model", model, "--prompt", prompt, "--max-tokens", maxTokens)

	// Create a pipe for the output of the command
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
