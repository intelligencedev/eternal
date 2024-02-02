# Tools

This folder contains a set of common tools and utilities used by multiple workflows.

## compress.go

Automates the process of compressing files and folders within a project. It is particularly useful for preparing static files for embedding into a Go binary, reducing the overall size of the compiled application.

### Functionality

- **File Compression**: The script traverses a specified directory and compresses all files within it using the GZIP format. This is particularly beneficial for compressing text-based files like HTML, CSS, and JavaScript, which often see significant size reductions when compressed.
- **Seamless Build Integration**: Integrated into the build process via the `go:generate` directive, this script ensures that all required files are compressed each time the project is built, keeping the embedded resources up-to-date and optimized for size.

### Usage

Use the go:generate directive in your Go files to invoke this script during the build process. For example:

```
// main.go

//go:generate go run build/compress.go
package main

// ... rest of your code ...
```

Run `go generate` in your project directory to execute the script and compress the files. The compressed files (`.gz` format) will be created in the same location as the original files.