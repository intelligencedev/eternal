package tools

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCompressFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create a test file
	testFileName := "test.txt"
	testContent := []byte("This is a test")
	err := afero.WriteFile(fs, testFileName, testContent, 0644)
	assert.NoError(t, err)

	// Compress the test file
	compressedFileName := testFileName + ".gz"
	err = CompressFile(fs, testFileName, compressedFileName)
	assert.NoError(t, err)

	// Read and decompress the file
	reader, err := fs.Open(compressedFileName)
	assert.NoError(t, err)
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	assert.NoError(t, err)
	defer gzipReader.Close()

	decompressedContent, err := io.ReadAll(gzipReader)
	assert.NoError(t, err)

	// Assert the decompressed content is the same as the original content
	assert.Equal(t, testContent, decompressedContent)
}

func TestCompressDirectory(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	log.Printf("Current working directory: %s", wd)

	srcDir := "../public"

	fs := afero.NewOsFs()

	compressedFileName := "../public/public.tar.gz"
	err = CompressDirectory(fs, srcDir, compressedFileName)
	assert.NoError(t, err)

	//fs = afero.NewMemMapFs()

	// Verify the compressed file exists
	exists, err := afero.Exists(fs, compressedFileName)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Decompress the file
	reader, err := fs.Open(compressedFileName)
	assert.NoError(t, err)
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	assert.NoError(t, err)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Read the tar file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatalf("Error reading tar file: %s", err)
		}

		t.Logf("File: %s", header.Name)

		// Read the file contents
		fileContent, err := io.ReadAll(tarReader)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}

		t.Logf("File contents: %s", fileContent)
	}

}
