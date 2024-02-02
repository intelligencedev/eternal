<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" height="200px" srcset="./public/img/eternal.png">
    <img alt="logo" height="200px" src="./public/img/eternal.png">
  </picture>
</div>

# Eternal

Eternal is an experimental platform for machine learning workflows.

## Features

- **Local Model Workflows**: Eternal manages the installation and default configurations of local machine learning models in GGUF, MLX and ONNX format for fast GPU and CPU inference on most hardware configurations.
- **Web Interface**: Provides a simple and modern web based interface for a consistent frontend experience regardless of operating system or client device.
- **WebSocket Support**: Utilizes WebSockets for real-time communication between the ml compute servers and the control plane.
- **Configurable**: Includes a configuration system that reads from a `config.yml` file. Adding new models is simple to understand and does not require a new application release.

## Configuration

Rename the provided `.config.yml` file to `config.yml` and place it in the same path as the application binary. Modify the contents for your environment and use case.

## Building the Application

Run `go mod tidy` to pull required third party Go packages first.

Eternal depends on the legendary [llama.cpp](https://github.com/ggerganov/llama.cpp) client library for GPU/CPU model inference. To build Eternal with llama support:

1. `git submodule init`
2. `git submodule update`
3. Navigate to the `./pkg/llm/local/gguf` path and build the llama.cpp library according to the instructions in its README for your arch. For MacOS, simply run `make`.
4. Copy the resulting binaries into the `pkg/llm/local/bin` path. Include the `ggml-metal.metal` if on MacOS M-Series arch.
5. Build Eternal with `go build`.
6. Eternal will bootstrap the dependent binaries and models depending on the selected workflow.

To run Eternal, ensure you have Go installed and the required dependencies are fetched. You can then build the application with `go build` and run the resulting binary. Make sure the `config.yml` file is present and properly set up before starting the server.

## Dependencies

Eternal relies on several Go packages, such as:

- `gofiber/fiber/v2`: For the web server framework.
- `gofiber/websocket/v2`: To manage WebSocket connections.
- `gorm.io/gorm`: For ORM support with SQLite.
- `github.com/spf13/afero`: For filesystem abstraction.
- `github.com/pterm/pterm`: For pretty terminal outputs.
- `embed`: For embedding static files into the binary.
- `github.com/gofiber/template/html/v2`: For HTML template rendering.

## Database Migrations

Upon initialization, Eternal attempts to auto-migrate the database schema for the following entities:

- `LanguageModels`
- `SelectedModels`
- `Chat`
- `llm.Model`

## Static Files

Static files for the frontend are embedded in the binary and served through the `/public` endpoint. They can be accessed directly or via the filesystem middleware provided by Fiber.

## Disclaimer

This README is a high-level overview of the Eternal application. Detailed setup instructions and a complete list of features, dependencies, and configurations should be consulted in the actual application documentation.
