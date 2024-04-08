# Eternal Embeddings CLI Tool

This CLI tool provides functionality to generate embeddings for a specified input file and retrieve top N similar words or chunks for a given prompt using a pre-trained BERT model.

## Installation

Before you can run the CLI tool, ensure you have Go installed on your system. You can download and install Go from [golang.org](https://golang.org/dl/).

## Usage

The CLI tool supports two main commands: `generate` and `retrieve`.

### Generate Command

The `generate` command is used to generate embeddings for the specified input file. To use this command, you must specify the input file path using the `--input-file` flag.

```bash
go run main.go generate --input-file <input file>
```

### Retrieve Command

The `retrieve` command is used to retrieve top N similar words or chunks for the given prompt. You must specify the prompt using the `--prompt` flag and can optionally specify the number of top similar words or chunks to retrieve using the `--top-n` flag.

```bash
go run main.go retrieve --prompt <prompt> [--top-n <number>]
```

If the `--top-n` flag is not specified, it defaults to 5.

### Examples

#### Generate Embeddings

To generate embeddings for an input file named `input.txt`, run the following command:

```bash
go run main.go generate --input-file input.txt
```

#### Retrieve Top Similar Words or Chunks

To retrieve the top 5 similar words or chunks for the prompt "machine learning", run the following command:

```bash
go run main.go retrieve --prompt "machine learning"
```

To retrieve the top 10 similar words or chunks for the prompt "artificial intelligence", run the following command:

```bash
go run main.go retrieve --prompt "artificial intelligence" --top-n 10
```

## Configuration

The tool uses flags to configure the model path, model name, and the limit for the number of dimensions in the embedding vector.

- `--model-path`: The path to the model directory (default is ".eternal/models/HF/").- `--model-name`: The name of the model (default is "avsolatorio/GIST-small-Embedding-v0").- `--limit`: The limit for the number of dimensions in the embedding vector (default is 128).

## Troubleshooting

If you encounter any issues while using the CLI tool, ensure that:

- The model path and model name are correctly specified and accessible.- The input file exists and is readable.- You have the necessary permissions to read and write to the files and directories.

For further assistance, please check the error messages provided by the CLI tool or consult the Go documentation