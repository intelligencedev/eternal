# MLX to Go

This Golang project is designed to execute a CLI command and stream its output. It specifically interfaces with a machine learning model to generate responses based on a given prompt.

## Overview

The main function of this project is to execute a CLI command (`mlx_lm.generate`) with specified parameters such as the model to use, the prompt for the model, and the maximum number of tokens to generate. The command's output is streamed and displayed in real time.

## Usage

The `main` function sets up and executes the CLI command with predefined parameters:

- **CLI Command**: `mlx_lm.generate`
- **Model**: `mlx-community/Phi-3-mini-128k-instruct-8bit`
- **Prompt**: `write a quicksort in python.`
- **Max Tokens**: `4096`

These parameters can be modified directly in the source code to fit different requirements or models.

## Functionality

### StreamCLICommand

The `StreamCLICommand` function is responsible for:
- Constructing the command with the necessary arguments.
- Handling the command's output and errors through pipes.
- Streaming the output to the console while capturing and logging any errors.

### Error Handling

The program logs fatal errors and exits if it fails to execute the CLI command or encounters issues during execution.

## System Requirements

Python:
This program requires the Python `mlx-lm` framework for LLM inference on Apple M-Series hardware.

Go:
This program is written in Go and requires the Go compiler to build and run. It also depends on the external CLI command (`mlx_lm.generate`) being available and correctly configured on the system where this program is executed.

## Installation

```
# with pip
$ pip install mlx-lm

# with conda
$ conda install -c conda-forge mlx-lm
```

To run this program:
1. Ensure you have Go installed on your machine.
2. Clone this repository.
3. Navigate to the directory containing the code.
4. Run `go build` to compile the program.
5. Execute the compiled program.

## Contributing

Contributions to this project are welcome. You can enhance it by:
- Adding parameter customization through command-line arguments.
- Improving error handling and logging.
- Extending compatibility with other CLI commands or models.

Please submit a pull request or open an issue if you have suggestions or improvements.

## License

Specify your license here or state that the project is unlicensed and available for free use.
