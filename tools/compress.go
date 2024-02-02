package tools

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func CompressFile(fs afero.Fs, src, dest string) error {
	reader, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	gw := gzip.NewWriter(writer)
	defer gw.Close()

	_, err = io.Copy(gw, reader)
	return err
}

func CompressDirectory(fs afero.Fs, srcDir, destFile string) error {
	// Create the output file
	outFile, err := fs.Create(destFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create the gzip writer
	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	// Create the tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Walk through the source directory
	return afero.Walk(fs, srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Open the file
		file, err := fs.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(filePath) // Ensure Unix-style paths

		// Write file header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Copy file content to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

		return nil
	})
}
