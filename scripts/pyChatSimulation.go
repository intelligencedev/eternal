// The provided Go program serves to set up a Python virtual environment,
// manage dependencies, and execute a Python script that simulates a chat
// by printing lorem ipsum sentences at random intervals.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Python script content
	pythonScript := `
import lorem
import time
import random

def simulate_chat():
	while True:
		message = lorem.sentence()
		print(f"Friend: {message}", flush=True)
		time_to_next_message = random.uniform(0.5, 3)
		time.sleep(time_to_next_message)
try:
    simulate_chat()
except KeyboardInterrupt:
    print("\\nChat simulation ended.", flush=True)
`

	// Write the Python script to a temporary file
	tmpFile, err := ioutil.TempFile("", "script-*.py")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(pythonScript)); err != nil {
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Bash script to set up the virtual environment and run the Python script
	bashScript := fmt.Sprintf(`
#!/bin/bash

# Check if the virtualenv exists
if [ ! -d "eternal" ]; then
    echo "Creating virtual environment..."
    python3 -m venv eternal
fi

echo "Activating virtual environment..."
source eternal/bin/activate

echo "Installing dependencies..."
pip install lorem

# Run the Python script
python3 %s
`, tmpFile.Name())

	// Run the bash script
	if err := RunShellScript(bashScript); err != nil {
		log.Fatalf("Failed to run script: %v", err)
	}
}

func RunShellScript(scriptContent string) error {
	cmd := exec.Command("bash", "-c", scriptContent)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Print(line)
	}

	return cmd.Wait()
}
